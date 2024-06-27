package imap

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chk-n/retry"
	"github.com/enuan/go-imap/parser"
	"github.com/rs/xid"
)

// AddSlashes adds slashes to double quotes
var AddSlashes = strings.NewReplacer(`"`, `\"`)

// RemoveSlashes removes slashes before double quotes
var RemoveSlashes = strings.NewReplacer(`\"`, `"`)

// Verbose outputs every command and its response with the IMAP server
var Verbose = false

// SkipResponses skips printing server responses in verbose mode
var SkipResponses = false

var lastResp string

// Dialer is basically an IMAP connection
type Dialer struct {
	conn      net.Conn
	Folder    string
	Username  string
	Password  string
	Host      string
	Port      int
	TLSConfig *tls.Config
	Connected bool
	ConnNum   int
}

var nextConnNum = 0
var nextConnNumMutex = sync.Mutex{}

func log(connNum int, folder string, msg interface{}) {
	var name string
	if len(folder) != 0 {
		name = fmt.Sprintf("IMAP%d:%s", connNum, folder)
	} else {
		name = fmt.Sprintf("IMAP%d", connNum)
	}
	fmt.Printf("%s %s: %s\n", time.Now().Format("2006-01-02 15:04:05.000000"), name, msg)
}

type Config struct {
	Username  string
	Password  string
	Host      string
	Port      int
	Secure    bool
	TLSConfig *tls.Config
}

// New makes a new imap
func New(cfg Config) (d *Dialer, err error) {
	nextConnNumMutex.Lock()
	connNum := nextConnNum
	nextConnNum++
	nextConnNumMutex.Unlock()

	r := retry.NewDefault()
	err = r.Do(func() (err error) {
		// on error, reconnect
		defer func() {
			if err != nil {
				d.Reconnect()
			}
		}()

		if Verbose {
			log(connNum, "", "establishing connection")
		}
		var conn net.Conn
		connStr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		if cfg.Secure {
			conn, err = tls.Dial("tcp", connStr, cfg.TLSConfig)
		} else {
			conn, err = net.Dial("tcp", connStr)
		}
		if err != nil {
			if Verbose {
				log(connNum, "", fmt.Sprintf("failed to connect: %s", err))
			}
			return err
		}
		d = &Dialer{
			conn:      conn,
			Username:  cfg.Username,
			Password:  cfg.Password,
			Host:      cfg.Host,
			Port:      cfg.Port,
			Connected: true,
			ConnNum:   connNum,
		}

		return d.Login(cfg.Username, cfg.Password)
	})
	if err != nil {
		if Verbose {
			log(connNum, "", "failed to establish connection")
		}
		if d != nil {
			d.Close()
		}
		return nil, err
	}

	return
}

// Clone returns a new connection with the same connection information
// as the one this is being called on
func (d *Dialer) Clone() (d2 *Dialer, err error) {
	d2, err = New(Config{
		Username:  d.Username,
		Password:  d.Password,
		Host:      d.Host,
		Port:      d.Port,
		TLSConfig: d.TLSConfig,
	})
	// d2.Verbose = d1.Verbose
	if d.Folder != "" {
		_, err = d2.SelectFolder(d.Folder)
		if err != nil {
			return nil, fmt.Errorf("imap clone: %s", err)
		}
	}
	return
}

// Close closes the imap connection
func (d *Dialer) Close() error {
	if !d.Connected {
		return nil
	}

	if Verbose {
		log(d.ConnNum, d.Folder, "closing connection")
	}
	d.Connected = false

	if d.conn == nil {
		return nil
	}

	if err := d.conn.Close(); err != nil {
		return fmt.Errorf("imap close: %s", err)
	}

	return nil
}

// Reconnect closes the current connection (if any) and establishes a new one
func (d *Dialer) Reconnect() (err error) {
	d.Close()
	if Verbose {
		log(d.ConnNum, d.Folder, "reopening connection")
	}
	d2, err := d.Clone()
	if err != nil {
		return fmt.Errorf("imap reconnect: %s", err)
	}
	*d = *d2
	return
}

const nl = "\r\n"

func dropNl(b []byte) []byte {
	if len(b) >= 1 && b[len(b)-1] == '\n' {
		if len(b) >= 2 && b[len(b)-2] == '\r' {
			return b[:len(b)-2]
		} else {
			return b[:len(b)-1]
		}
	}
	return b
}

var atom = regexp.MustCompile(`{\d+}$`)

// Exec executes the command on the imap connection
func (d *Dialer) Exec(command string, buildResponse bool, processLine func(line []byte) error) (response string, err error) {
	var resp strings.Builder
	r := retry.NewDefault()
	err = r.Do(func() (err error) {
		// on error, reconnect
		defer func() {
			if err != nil {
				d.Reconnect()
			}
		}()

		tag := []byte(fmt.Sprintf("%X", xid.New()))

		c := fmt.Sprintf("%s %s\r\n", tag, command)

		if Verbose {
			log(d.ConnNum, d.Folder, strings.Replace(fmt.Sprintf("%s %s", "->", strings.TrimSpace(c)), fmt.Sprintf(`"%s"`, d.Password), `"****"`, -1))
		}

		_, err = d.conn.Write([]byte(c))
		if err != nil {
			return err
		}

		r := bufio.NewReader(d.conn)

		if buildResponse {
			resp = strings.Builder{}
		}
		var line []byte
		for err == nil {
			line, err = r.ReadBytes('\n')
			for {
				if a := atom.Find(dropNl(line)); a != nil {
					// fmt.Printf("%s\n", a)
					var n int
					n, err = strconv.Atoi(string(a[1 : len(a)-1]))
					if err != nil {
						return err
					}

					buf := make([]byte, n)
					_, err = io.ReadFull(r, buf)
					if err != nil {
						return err
					}
					line = append(line, buf...)

					buf, err = r.ReadBytes('\n')
					if err != nil {
						return err
					}
					line = append(line, buf...)

					continue
				}
				break
			}

			if Verbose && !SkipResponses {
				log(d.ConnNum, d.Folder, fmt.Sprintf("<- %s", dropNl(line)))
			}

			// if strings.Contains(string(line), "--00000000000030095105741e7f1f") {
			// 	f, _ := ioutil.TempFile("", "")
			// 	ioutil.WriteFile(f.Name(), line, 0777)
			// 	fmt.Println(f.Name())
			// }

			// XID project is returning 40-byte tags. The code was originally hardcoded 16 digits.
			taglen := len(tag)
			oklen := 3
			if len(line) >= taglen+oklen && bytes.Equal(line[:taglen], tag) {
				if !bytes.Equal(line[taglen+1:taglen+oklen], []byte("OK")) {
					err = fmt.Errorf("imap command failed: %s", line[taglen+oklen+1:])
					return err
				}
				break
			}

			if processLine != nil {
				if err = processLine(line); err != nil {
					return err
				}
			}
			if buildResponse {
				resp.Write(line)
			}
		}
		return nil
	})
	if err != nil {
		if Verbose {
			log(d.ConnNum, d.Folder, "All retries failed")
		}
		return "", err
	}

	if buildResponse {
		if resp.Len() != 0 {
			lastResp = resp.String()
			return lastResp, nil
		}
		return "", nil
	}
	return
}

// Login attempts to login
func (d *Dialer) Login(username string, password string) (err error) {
	_, err = d.Exec(fmt.Sprintf(`LOGIN "%s" "%s"`, AddSlashes.Replace(username), AddSlashes.Replace(password)), false, nil)
	return
}

// SelectFolder selects a folder
func (d *Dialer) SelectFolder(folder string) (uidValidity uint32, err error) {
	resp, err := d.Exec(`EXAMINE "`+AddSlashes.Replace(folder)+`"`, true, nil)
	if err != nil {
		return
	}
	uidValidity, err = parser.ParseExamineResponse(resp)
	d.Folder = folder
	return
}

// Move a read email to a specified folder
func (d *Dialer) MoveEmail(uid int, folder string) (err error) {
	_, err = d.Exec(`UID MOVE `+strconv.Itoa(uid)+` "`+AddSlashes.Replace(folder)+`"`, true, nil)
	if err != nil {
		return
	}
	d.Folder = folder
	return nil
}

// GetUIDs returns the UIDs in the current folder that match the search
func (d *Dialer) GetUIDs(search string) (uids []uint32, err error) {
	uids = make([]uint32, 0)
	r, err := d.Exec(`UID SEARCH `+search, true, nil)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(r, nl)
	for _, l := range lines {
		prefix := "* SEARCH"
		if strings.HasPrefix(l, prefix) {
			suffix := l[len(prefix):]
			for _, atom := range strings.Fields(suffix) {
				uid, err := strconv.ParseUint(atom, 10, 32)
				if err != nil {
					return nil, err
				}
				uids = append(uids, uint32(uid))
			}
		}
	}
	return uids, nil
}

func (d *Dialer) GetEmailByUID(uid uint32) (string, error) {
	cmd := fmt.Sprintf("UID FETCH %d BODY[]", uid)
	r, err := d.Exec(cmd, true, nil)
	if err != nil {
		return "", err
	}
	fetchInfo, err := parser.ParseFetchResponse(r)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(fetchInfo["BODY[]"], nl), nil
}
