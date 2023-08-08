package domain

import (
	"fyne.io/fyne/v2/data/binding"
)

// UpsStatusValueBindings bond datapoints for UPS status display
type UpsStatusValueBindings struct {
	Host       *ApcHost
	Bselftest  binding.String
	Bnumxfers  binding.String
	Blastxfer  binding.String
	Bxonbatt   binding.String
	Bxoffbatt  binding.String
	Btonbatt   binding.String
	Bcumonbatt binding.String
	Bhostname  binding.String
	Bupsname   binding.String
	Bmaster    binding.String
	Blinev     binding.String
	Bbattv     binding.String
	Bbcharge   binding.String
	Bloadpct   binding.String
	Btimeleft  binding.String
	Bversion   binding.String
	Bcable     binding.String
	Bdriver    binding.String
	Bupsmode   binding.String
	Bstarttime binding.String
	Bstatus    binding.String
	Bmodel     binding.String
	Bserialno  binding.String
	Bmandate   binding.String
	Bfirmware  binding.String
	Bbattdate  binding.String
	Bitemp     binding.String
}

// NewUpsStatusValueBindings creates a new collection of string to be bound
// to the view elements for UPS.
func NewUpsStatusValueBindings(h *ApcHost) *UpsStatusValueBindings {
	return &UpsStatusValueBindings{
		Host:       h,
		Bselftest:  binding.NewString(),
		Bnumxfers:  binding.NewString(),
		Blastxfer:  binding.NewString(),
		Bxonbatt:   binding.NewString(),
		Bxoffbatt:  binding.NewString(),
		Btonbatt:   binding.NewString(),
		Bcumonbatt: binding.NewString(),
		Bhostname:  binding.NewString(),
		Bupsname:   binding.NewString(),
		Bmaster:    binding.NewString(),
		Blinev:     binding.NewString(),
		Bbattv:     binding.NewString(),
		Bbcharge:   binding.NewString(),
		Bloadpct:   binding.NewString(),
		Btimeleft:  binding.NewString(),
		Bversion:   binding.NewString(),
		Bcable:     binding.NewString(),
		Bdriver:    binding.NewString(),
		Bupsmode:   binding.NewString(),
		Bstarttime: binding.NewString(),
		Bstatus:    binding.NewString(),
		Bmodel:     binding.NewString(),
		Bserialno:  binding.NewString(),
		Bmandate:   binding.NewString(),
		Bfirmware:  binding.NewString(),
		Bbattdate:  binding.NewString(),
		Bitemp:     binding.NewString(),
	}
}

// Apply reads slices from provided channel and assigns value to bond strings
func (b *UpsStatusValueBindings) Apply(params map[string]string) {
	h := b.Host
	for k, v := range params {
		switch k {
		case "SELFTEST":
			_ = b.Bselftest.Set(v)
		case "NUMXFERS":
			_ = b.Bnumxfers.Set(v)
			_ = h.Bnumxfers.Set(v)
		case "LASTXFER":
			_ = b.Blastxfer.Set(v)
		case "XONBATT":
			_ = b.Bxonbatt.Set(v)
			_ = h.Bxonbatt.Set(v)
		case "XOFFBATT":
			_ = b.Bxoffbatt.Set(v)
			_ = h.Bxoffbatt.Set(v)
		case "TONBATT":
			_ = b.Btonbatt.Set(v)
		case "CUMONBATT":
			_ = b.Bcumonbatt.Set(v)
			_ = h.Bcumonbatt.Set(v)
		case "HOSTNAME":
			_ = b.Bhostname.Set(v)
		case "UPSNAME":
			_ = b.Bupsname.Set(v)
		case "MASTER":
			_ = b.Bmaster.Set(v)
		case "LINEV":
			_ = b.Blinev.Set(v)
			_ = h.Blinev.Set(v)
		case "BATTV":
			_ = b.Bbattv.Set(v)
		case "BCHARGE":
			_ = b.Bbcharge.Set(v)
			_ = h.Bbcharge.Set(v)
		case "LOADPCT":
			_ = b.Bloadpct.Set(v)
			_ = h.Bloadpct.Set(v)
		case "TIMELEFT":
			_ = b.Btimeleft.Set(v)
		case "VERSION":
			_ = b.Bversion.Set(v)
		case "CABLE":
			_ = b.Bcable.Set(v)
		case "DRIVER":
			_ = b.Bdriver.Set(v)
		case "UPSMODE":
			_ = b.Bupsmode.Set(v)
		case "STARTTIME":
			_ = b.Bstarttime.Set(v)
		case "STATUS":
			_ = b.Bstatus.Set(v)
		case "MODEL":
			_ = b.Bmodel.Set(v)
		case "SERIALNO":
			_ = b.Bserialno.Set(v)
		case "MANDATE":
			_ = b.Bmandate.Set(v)
		case "FIRMWARE":
			_ = b.Bfirmware.Set(v)
		case "BATTDATE":
			_ = b.Bbattdate.Set(v)
		case "ITEMP":
			_ = b.Bitemp.Set(v)
		}
	}
}
func (b *UpsStatusValueBindings) UnbindUpsData() {
	/*
	 * No way to unBind values
	 */
}
func (b *UpsStatusValueBindings) IsNil() bool {
	return (b == nil)
}
