package main

import (
	"fmt"

	"github.com/enuan/go-imap/parser"
)

var s = "* 160 FETCH (BODY[] {3770}\r\nReceived: from mail.tns.it ([89.96.139.115]) by tnsmx02mil.uvet.com with Microsoft SMTPSVC(6.0.3790.4675);\r\n\t Wed, 10 Jul 2024 16:05:14 +0200\r\nReceived: from mail-ed1-f52.google.com ([209.85.208.52]) by mail.tns.it over TLS secured channel\r\n with Microsoft SMTPSVC(10.0.14393.4169); Wed, 10 Jul 2024 16:05:11 +0200\r\nReceived: by mail-ed1-f52.google.com with SMTP id 4fb4d7f45d1cf-58ba3e38027so7935136a12.1\r\n        for <test-enuan@uvetgbt.com>; Wed, 10 Jul 2024 07:05:11 -0700 (PDT)\r\nDKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;        d=enuan-com.20230601.gappssmtp.com;\r\n s=20230601; t=1720620310; x=1721225110; darn=uvetgbt.com;        h=to:subject:message-id:date:from:mime-version:from:to:cc:subject\r\n         :date:message-id:reply-to;        bh=87HBLqfWZQbJ77DHIcRnl2DIjpKG4+QZ1TaxMIwKTcM=;\r\n        b=sxtA913FiYf9OqKt42vi0fOqQ8VfX6Wecx/BdcpOWdyHY55FN/aU2c0Q1VVzkNLreu   \r\n      YabTfzHE6OIjGEmpVIg1KAdpqfSMIF0V1p4noL10wasLC+G4PFd6I11BMsjPLD+enbP5     \r\n    quwuMCFTDvwlldAwcMjoS0VI5LHLm5lzmmNhPvTAskEcXXxAhwH29wPPxJPd52cpcS/e       \r\n  Zfq5NS0eU0ZBSvwrhlijRdAwIdYUjezNcDnzHSSBpJTazDRITRFVabI/2cfxGKyVNiIu         6pDuL946CYncJMkuvu45gKAmgSgcawa/fppTVBEyHC4rjjCqNOAEOupu92eZR9LwV/OK         M1/A==\r\nX-Google-DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;        d=1e100.net; s=20230601; t=1720620310;\r\n x=1721225110;        h=to:subject:message-id:date:from:mime-version:x-gm-message-state\r\n         :from:to:cc:subject:date:message-id:reply-to;        bh=87HBLqfWZQbJ77DHIcRnl2DIjpKG4+QZ1TaxMIwKTcM=;\r\n        b=PextAIOfd6c1MSXcwABRq3xo2Cj6YLIJgQ57IWZjlLcGHPIqVmw2Crk/ms7iP35xJT   \r\n      EOToPXFLFGkWlKHdbDb47ojb/mVugEmNKcJjue3eNxda0GlYPM0wmfi72i6vJd3Z/Qqo     \r\n    nhTVgt/u10ZN1gQQSeuXWFaDcEKx0BoeMzvIswVq8YMo8vJazf1iljul3MeOdTsbBDE1       \r\n  uKLsC2PogzWy7xlIZ33G0Y4XYhjsVvXvKs90Iltioc+cvsw0kkemHDFKXUvTfDKGCCU4         wsSo9l8CAGQqZQnHF2pRNbHFuLWnD3qD6RzQkOcPL5xsoTMtE5vGkKzxA4Tlbb6ssFQp         1RxA==\r\nX-Gm-Message-State: AOJu0YxA7Xm3D7exfSC47XbgorzLm5nbiKqT7ccGTz0lExzIRv60UFFPgU6YaP/s8nXaJSx46zQ8GM82QzfrFQVDKK7mBvAHoFTUunMPxwq0F5Q73THnoudfi+H1TRJUJilUkT8oFiL6Si2DjYek1ZmPEpBusj4YNRXz67K2Mv78EPQ8sac=\r\nX-Google-Smtp-Source: AGHT+IG1MRxtmmlIDTPw7D/UqWNZL3GSXfYhEiRMQhhTDo/WpaaP731vZTRVLe/M9ZAg8wkT5IjOz7RXjaaDjoOsBzY=\r\nX-Received: by 2002:a05:6402:27ce:b0:57c:80bf:9267 with SMTP id 4fb4d7f45d1cf-594baa8bc33mr5124832a12.6.1720620309900; Wed, 10 Jul 2024 07:05:09 -0700 (PDT)\r\nMIME-Version: 1.0\r\nFrom: Paolo Losi <paolo@enuan.com>\r\nDate: Wed, 10 Jul 2024 16:04:58 +0200\r\nMessage-ID: <CAP=2L=HJpbMLiHzTh+a4ANjGmuFsD6_Vc0ZQ1QGaXzcoY-1XOw@mail.gmail.com>\r\nSubject: test move 9\r\nTo: test-enuan@uvetgbt.com\r\nContent-Type: multipart/alternative; boundary=\"0000000000005d8c42061ce5224a\"\r\nReturn-Path: paolo@enuan.com\r\nX-OriginalArrivalTime: 10 Jul 2024 14:05:11.0326 (UTC) FILETIME=[2DBE43E0:01DAD2D2]\r\n\r\nThis is a multi-part message in MIME format.\r\n\r\n--0000000000005d8c42061ce5224a\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\nX-EC0D2A8E-5CB7-4969-9C36-46D859D137BE-PartID: {16CAC69C-5616-4917-8C9C-F21CF9A5ADE4}\r\n\r\n-- \r\nPaolo Losi\r\nmobile:   +39 348 7705261\r\nSoftware Engineering @ www.enuan.com\r\n\r\n\r\n\r\nENUAN Srl\r\nVia XX Settembre, 12 - 29121 Piacenza\r\n\r\n--0000000000005d8c42061ce5224a\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nContent-Transfer-Encoding: quoted-printable\r\nX-EC0D2A8E-5CB7-4969-9C36-46D859D137BE-PartID: {6C1C760F-A1CA-499D-8D18-6687EEE1BA28}\r\n\r\n<div><br><div><br></div><span class=3D\"gmail_signature_prefix\">-- </span><br=\r\n><div class=3D\"gmail_signature\"><div>Paolo Losi<br>mobile:=C2=A0=C2=A0 +39 3=\r\n48 7705261<div>Software Engineering @ <a href=3D\"http://www.enuan.com\">www.e=\r\nnuan.com</a><br><br><br><br>ENUAN Srl<br>Via XX Settembre, 12 - 29121 Piacen=\r\nza</div></div></div></div>\r\n\r\n--0000000000005d8c42061ce5224a--\r\n UID 163 FLAGS (\\Seen \\Recent))\r\n"

func main() {
	m, err := parser.ParseFetchResponse(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}