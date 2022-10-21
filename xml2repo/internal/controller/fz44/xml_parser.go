package fz44

import (
	"encoding/xml"
	"io"
	"os"
	"strings"
	"time"

	"github.com/b2b2b-pro/lib/object"
	"github.com/b2b2b-pro/lib/repo_client"
	"go.uber.org/zap"
)

type Parser struct {
	repo    *repo_client.RepoGRPC
	decoder *xml.Decoder
}

func NewParser(db *repo_client.RepoGRPC) *Parser {
	zap.S().Debug("Configuring XML Parser.")
	xp := &Parser{}
	xp.repo = db
	return xp
}

func (xp *Parser) parseOrigin(tc xml.StartElement) object.Origin {
	var err error
	orn := object.Origin{}
	err = xp.decoder.DecodeElement(&orn.Description, &tc)
	if err != nil {
		zap.S().Debugf("decoder.DecodeElement(org.Description, &tc) error: %v\n", err)
	}
	zap.S().Debugf("ссылка на госзакупки: %v\n", orn.Description)
	// создать origin, получить его ID и использовать для obligation
	orn.ID, err = xp.repo.CreateOrigin(orn)
	if err != nil {
		zap.S().Debugf("xp.repo.CreateOrigin(orn) error: %v", err)
	}
	return orn
}

func (xp *Parser) parseEntity(end string) object.Entity {
	var en object.Entity
	var err error
fl:
	for {
		var d xml.Token
		d, err = xp.decoder.Token()
		if err != nil {
			zap.S().Errorf("tok, err := decoder.Token(), err: %v\n", err)
		}
		if d == nil {
			break
		}
		switch td := d.(type) {
		case xml.StartElement:
			tdname := strings.ToLower(td.Name.Local)
			switch tdname {
			case "inn":
				xp.decoder.DecodeElement(&en.INN, &td)
			case "kpp":
				xp.decoder.DecodeElement(&en.KPP, &td)
			case "shortname":
				xp.decoder.DecodeElement(&en.ShortName, &td)
			case "fullname":
				xp.decoder.DecodeElement(&en.FullName, &td)
			}
		case xml.EndElement:
			if td.Name.Local == end {
				zap.S().Debugf("entity: %v\n", en)
				break fl
			}
		}
	}
	// создать entity
	en.ID, err = xp.repo.CreateEntity(en)
	if err != nil {
		zap.S().Errorf("repo.CreateEntity error: %v\n", err)
	}
	return en
}

func (xp *Parser) parsePayment(end string) object.Obligation {
	var ob object.Obligation
	var err error
	var m, y int
fl:
	for {
		var d xml.Token
		d, err = xp.decoder.Token()
		if err != nil {
			zap.S().Errorf("tok, err := decoder.Token(), err: %v\n", err)
		}
		if d == nil {
			break
		}
		switch td := d.(type) {
		case xml.StartElement:
			tdname := strings.ToLower(td.Name.Local)
			switch tdname {
			case "paymentmonth":
				xp.decoder.DecodeElement(&m, &td)
			case "paymentyear":
				xp.decoder.DecodeElement(&y, &td)
			case "paymentsumrur":
				xp.decoder.DecodeElement(&ob.Cost, &td)
			}
		case xml.EndElement:
			if td.Name.Local == end {
				zap.S().Debugf("obligation: %v\n", ob)
				break fl
			}
		}
	}
	// вычислить payment_date
	ob.Date = object.NewPaymentDate(time.Date(y, time.Month(m+1), 0, 0, 0, 0, 0, time.UTC))
	return ob
}

func (xp *Parser) ParseXML(fname string) error {
	var err error
	f, err := os.Open(fname)
	if err != nil {
		zap.S().Debug("ParseXML can't open file ", fname, ", error: ", err, "\n")
		return err
	}
	defer f.Close()

	xp.decoder = xml.NewDecoder(f)
	var c xml.Token
	var orn object.Origin
	var csmr, splr object.Entity
	// рассчитываем, что исполнитель по контракту будет один, но могут быть несколько платежей
	// хотя, возможно, несколько платежей обусловлены технологией, поэтому может быть их не надо учитывать
	var obPay []object.Obligation
	var obPrc object.Obligation
	fp := false // один раз заполнить price
	fd := false // один раз заполнить enddate

	for {
		c, err = xp.decoder.Token()
		if err == io.EOF {
			zap.S().Debugf("parseXML %s EOF\n", fname)
			break
		}
		if err != nil {
			zap.S().Errorf("parseXML %s error: %v\n", fname, err)
			return err
		}
		switch tc := c.(type) {
		case xml.StartElement:
			tcname := strings.ToLower(tc.Name.Local)
			switch {
			case strings.HasPrefix(tcname, "customer"):
				end := tc.Name.Local
				csmr = xp.parseEntity(end)
				zap.S().Debugf("Customer: %v\n", csmr)

			case strings.HasPrefix(tcname, "supplier"):
				end := tc.Name.Local
				splr = xp.parseEntity(end)
				zap.S().Debugf("Supplier: %v\n", splr)

			case strings.HasPrefix(tcname, "href"):
				orn = xp.parseOrigin(tc)
				zap.S().Debugf("Origin: %v\n", orn)

			case strings.HasPrefix(tcname, "payments"):
				end := tc.Name.Local
				ob := xp.parsePayment(end)
				obPay = append(obPay, ob)
				zap.S().Debugf("Payment: %v\n", ob)

			case strings.HasPrefix(tcname, "pricerur"):
				if fp {
					continue
				}
				err = xp.decoder.DecodeElement(&obPrc.Cost, &tc)
				if err != nil {
					zap.S().Debugf("decoder.DecodeElement pricerur error: %v\n", err)
				} else {
					fp = true
				}
				zap.S().Debugf("Price: %v\n", obPrc)

			case strings.HasPrefix(tcname, "enddate"):
				if fd {
					continue
				}
				err = xp.decoder.DecodeElement(&obPrc.Date, &tc)
				if err != nil {
					zap.S().Debugf("decoder.DecodeElement enddate error: %v\n", err)
				} else {
					fd = true
				}
			}
		}
	}
	if err != io.EOF && err != nil {
		zap.S().Debugf("разбор XML завершился с ошибкой: %v\n", err)
		return err
	}
	for _, x := range obPay {
		x.CreditorID = splr.ID
		x.DebtorID = csmr.ID
		x.OriginID = orn.ID
		xp.repo.CreateObligation(x)
	}

	// TODO перенести в конец
	obPrc.CreditorID = splr.ID
	obPrc.DebtorID = csmr.ID
	obPrc.OriginID = orn.ID

	zap.S().Debugf("информация о долге, полученная из Price: %v\n", obPrc)

	if len(obPay) > 0 {
		return nil
	}

	xp.repo.CreateObligation(obPrc)

	return nil
}
