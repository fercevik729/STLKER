// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/watcher.proto

package protos

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Currencies is the enum which represents the allowed currencies for the API
// USD is the default currency
type Currencies int32

const (
	Currencies_USD Currencies = 0
	Currencies_EUR Currencies = 1
	Currencies_JPY Currencies = 2
	Currencies_BGN Currencies = 3
	Currencies_CZK Currencies = 4
	Currencies_DKK Currencies = 5
	Currencies_GBP Currencies = 6
	Currencies_HUF Currencies = 7
	Currencies_PLN Currencies = 8
	Currencies_RON Currencies = 9
	Currencies_SEK Currencies = 10
	Currencies_CHF Currencies = 11
	Currencies_ISK Currencies = 12
	Currencies_NOK Currencies = 13
	Currencies_HRK Currencies = 14
	Currencies_RUB Currencies = 15
	Currencies_TRY Currencies = 16
	Currencies_AUD Currencies = 17
	Currencies_BRL Currencies = 18
	Currencies_CAD Currencies = 19
	Currencies_CNY Currencies = 20
	Currencies_HKD Currencies = 21
	Currencies_IDR Currencies = 22
	Currencies_ILS Currencies = 23
	Currencies_INR Currencies = 24
	Currencies_KRW Currencies = 25
	Currencies_MXN Currencies = 26
	Currencies_MYR Currencies = 27
	Currencies_NZD Currencies = 28
	Currencies_PHP Currencies = 29
	Currencies_SGD Currencies = 30
	Currencies_THB Currencies = 31
	Currencies_ZAR Currencies = 32
)

var Currencies_name = map[int32]string{
	0:  "USD",
	1:  "EUR",
	2:  "JPY",
	3:  "BGN",
	4:  "CZK",
	5:  "DKK",
	6:  "GBP",
	7:  "HUF",
	8:  "PLN",
	9:  "RON",
	10: "SEK",
	11: "CHF",
	12: "ISK",
	13: "NOK",
	14: "HRK",
	15: "RUB",
	16: "TRY",
	17: "AUD",
	18: "BRL",
	19: "CAD",
	20: "CNY",
	21: "HKD",
	22: "IDR",
	23: "ILS",
	24: "INR",
	25: "KRW",
	26: "MXN",
	27: "MYR",
	28: "NZD",
	29: "PHP",
	30: "SGD",
	31: "THB",
	32: "ZAR",
}

var Currencies_value = map[string]int32{
	"USD": 0,
	"EUR": 1,
	"JPY": 2,
	"BGN": 3,
	"CZK": 4,
	"DKK": 5,
	"GBP": 6,
	"HUF": 7,
	"PLN": 8,
	"RON": 9,
	"SEK": 10,
	"CHF": 11,
	"ISK": 12,
	"NOK": 13,
	"HRK": 14,
	"RUB": 15,
	"TRY": 16,
	"AUD": 17,
	"BRL": 18,
	"CAD": 19,
	"CNY": 20,
	"HKD": 21,
	"IDR": 22,
	"ILS": 23,
	"INR": 24,
	"KRW": 25,
	"MXN": 26,
	"MYR": 27,
	"NZD": 28,
	"PHP": 29,
	"SGD": 30,
	"THB": 31,
	"ZAR": 32,
}

func (x Currencies) String() string {
	return proto.EnumName(Currencies_name, int32(x))
}

func (Currencies) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_606cf1056908c0d2, []int{0}
}

type TickerRequest struct {
	Ticker               string     `protobuf:"bytes,1,opt,name=Ticker,proto3" json:"Ticker,omitempty"`
	Destination          Currencies `protobuf:"varint,2,opt,name=Destination,proto3,enum=Currencies" json:"Destination,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *TickerRequest) Reset()         { *m = TickerRequest{} }
func (m *TickerRequest) String() string { return proto.CompactTextString(m) }
func (*TickerRequest) ProtoMessage()    {}
func (*TickerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_606cf1056908c0d2, []int{0}
}

func (m *TickerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TickerRequest.Unmarshal(m, b)
}
func (m *TickerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TickerRequest.Marshal(b, m, deterministic)
}
func (m *TickerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TickerRequest.Merge(m, src)
}
func (m *TickerRequest) XXX_Size() int {
	return xxx_messageInfo_TickerRequest.Size(m)
}
func (m *TickerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TickerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TickerRequest proto.InternalMessageInfo

func (m *TickerRequest) GetTicker() string {
	if m != nil {
		return m.Ticker
	}
	return ""
}

func (m *TickerRequest) GetDestination() Currencies {
	if m != nil {
		return m.Destination
	}
	return Currencies_USD
}

type TickerResponse struct {
	Symbol               string   `protobuf:"bytes,1,opt,name=Symbol,proto3" json:"Symbol,omitempty"`
	Open                 string   `protobuf:"bytes,2,opt,name=Open,proto3" json:"Open,omitempty"`
	High                 string   `protobuf:"bytes,3,opt,name=High,proto3" json:"High,omitempty"`
	Low                  string   `protobuf:"bytes,4,opt,name=Low,proto3" json:"Low,omitempty"`
	Price                string   `protobuf:"bytes,5,opt,name=Price,proto3" json:"Price,omitempty"`
	Destination          string   `protobuf:"bytes,6,opt,name=Destination,proto3" json:"Destination,omitempty"`
	Volume               string   `protobuf:"bytes,7,opt,name=Volume,proto3" json:"Volume,omitempty"`
	LTD                  string   `protobuf:"bytes,8,opt,name=LTD,proto3" json:"LTD,omitempty"`
	PrevClose            string   `protobuf:"bytes,9,opt,name=PrevClose,proto3" json:"PrevClose,omitempty"`
	Change               string   `protobuf:"bytes,10,opt,name=Change,proto3" json:"Change,omitempty"`
	PercentChange        string   `protobuf:"bytes,11,opt,name=PercentChange,proto3" json:"PercentChange,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TickerResponse) Reset()         { *m = TickerResponse{} }
func (m *TickerResponse) String() string { return proto.CompactTextString(m) }
func (*TickerResponse) ProtoMessage()    {}
func (*TickerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_606cf1056908c0d2, []int{1}
}

func (m *TickerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TickerResponse.Unmarshal(m, b)
}
func (m *TickerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TickerResponse.Marshal(b, m, deterministic)
}
func (m *TickerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TickerResponse.Merge(m, src)
}
func (m *TickerResponse) XXX_Size() int {
	return xxx_messageInfo_TickerResponse.Size(m)
}
func (m *TickerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TickerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TickerResponse proto.InternalMessageInfo

func (m *TickerResponse) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *TickerResponse) GetOpen() string {
	if m != nil {
		return m.Open
	}
	return ""
}

func (m *TickerResponse) GetHigh() string {
	if m != nil {
		return m.High
	}
	return ""
}

func (m *TickerResponse) GetLow() string {
	if m != nil {
		return m.Low
	}
	return ""
}

func (m *TickerResponse) GetPrice() string {
	if m != nil {
		return m.Price
	}
	return ""
}

func (m *TickerResponse) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *TickerResponse) GetVolume() string {
	if m != nil {
		return m.Volume
	}
	return ""
}

func (m *TickerResponse) GetLTD() string {
	if m != nil {
		return m.LTD
	}
	return ""
}

func (m *TickerResponse) GetPrevClose() string {
	if m != nil {
		return m.PrevClose
	}
	return ""
}

func (m *TickerResponse) GetChange() string {
	if m != nil {
		return m.Change
	}
	return ""
}

func (m *TickerResponse) GetPercentChange() string {
	if m != nil {
		return m.PercentChange
	}
	return ""
}

type CompanyResponse struct {
	Ticker               string   `protobuf:"bytes,1,opt,name=Ticker,proto3" json:"Ticker,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Exchange             string   `protobuf:"bytes,3,opt,name=Exchange,proto3" json:"Exchange,omitempty"`
	Sector               string   `protobuf:"bytes,4,opt,name=Sector,proto3" json:"Sector,omitempty"`
	MarketCap            string   `protobuf:"bytes,5,opt,name=MarketCap,proto3" json:"MarketCap,omitempty"`
	PERatio              string   `protobuf:"bytes,6,opt,name=PERatio,proto3" json:"PERatio,omitempty"`
	PEGRatio             string   `protobuf:"bytes,7,opt,name=PEGRatio,proto3" json:"PEGRatio,omitempty"`
	DivPerShare          string   `protobuf:"bytes,8,opt,name=DivPerShare,proto3" json:"DivPerShare,omitempty"`
	EPS                  string   `protobuf:"bytes,9,opt,name=EPS,proto3" json:"EPS,omitempty"`
	RevPerShare          string   `protobuf:"bytes,10,opt,name=RevPerShare,proto3" json:"RevPerShare,omitempty"`
	ProfitMargin         string   `protobuf:"bytes,11,opt,name=ProfitMargin,proto3" json:"ProfitMargin,omitempty"`
	YearHigh             string   `protobuf:"bytes,12,opt,name=YearHigh,proto3" json:"YearHigh,omitempty"`
	YearLow              string   `protobuf:"bytes,13,opt,name=YearLow,proto3" json:"YearLow,omitempty"`
	SharesOutstanding    string   `protobuf:"bytes,14,opt,name=SharesOutstanding,proto3" json:"SharesOutstanding,omitempty"`
	PriceToBookRatio     string   `protobuf:"bytes,15,opt,name=PriceToBookRatio,proto3" json:"PriceToBookRatio,omitempty"`
	Beta                 string   `protobuf:"bytes,16,opt,name=Beta,proto3" json:"Beta,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CompanyResponse) Reset()         { *m = CompanyResponse{} }
func (m *CompanyResponse) String() string { return proto.CompactTextString(m) }
func (*CompanyResponse) ProtoMessage()    {}
func (*CompanyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_606cf1056908c0d2, []int{2}
}

func (m *CompanyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompanyResponse.Unmarshal(m, b)
}
func (m *CompanyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompanyResponse.Marshal(b, m, deterministic)
}
func (m *CompanyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompanyResponse.Merge(m, src)
}
func (m *CompanyResponse) XXX_Size() int {
	return xxx_messageInfo_CompanyResponse.Size(m)
}
func (m *CompanyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CompanyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CompanyResponse proto.InternalMessageInfo

func (m *CompanyResponse) GetTicker() string {
	if m != nil {
		return m.Ticker
	}
	return ""
}

func (m *CompanyResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CompanyResponse) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *CompanyResponse) GetSector() string {
	if m != nil {
		return m.Sector
	}
	return ""
}

func (m *CompanyResponse) GetMarketCap() string {
	if m != nil {
		return m.MarketCap
	}
	return ""
}

func (m *CompanyResponse) GetPERatio() string {
	if m != nil {
		return m.PERatio
	}
	return ""
}

func (m *CompanyResponse) GetPEGRatio() string {
	if m != nil {
		return m.PEGRatio
	}
	return ""
}

func (m *CompanyResponse) GetDivPerShare() string {
	if m != nil {
		return m.DivPerShare
	}
	return ""
}

func (m *CompanyResponse) GetEPS() string {
	if m != nil {
		return m.EPS
	}
	return ""
}

func (m *CompanyResponse) GetRevPerShare() string {
	if m != nil {
		return m.RevPerShare
	}
	return ""
}

func (m *CompanyResponse) GetProfitMargin() string {
	if m != nil {
		return m.ProfitMargin
	}
	return ""
}

func (m *CompanyResponse) GetYearHigh() string {
	if m != nil {
		return m.YearHigh
	}
	return ""
}

func (m *CompanyResponse) GetYearLow() string {
	if m != nil {
		return m.YearLow
	}
	return ""
}

func (m *CompanyResponse) GetSharesOutstanding() string {
	if m != nil {
		return m.SharesOutstanding
	}
	return ""
}

func (m *CompanyResponse) GetPriceToBookRatio() string {
	if m != nil {
		return m.PriceToBookRatio
	}
	return ""
}

func (m *CompanyResponse) GetBeta() string {
	if m != nil {
		return m.Beta
	}
	return ""
}

func init() {
	proto.RegisterEnum("Currencies", Currencies_name, Currencies_value)
	proto.RegisterType((*TickerRequest)(nil), "TickerRequest")
	proto.RegisterType((*TickerResponse)(nil), "TickerResponse")
	proto.RegisterType((*CompanyResponse)(nil), "CompanyResponse")
}

func init() { proto.RegisterFile("protos/watcher.proto", fileDescriptor_606cf1056908c0d2) }

var fileDescriptor_606cf1056908c0d2 = []byte{
	// 748 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x94, 0x5f, 0x73, 0xea, 0x44,
	0x18, 0xc6, 0xe5, 0x4f, 0xf9, 0xb3, 0x14, 0x78, 0xcf, 0x5a, 0x8f, 0xb1, 0x1e, 0x95, 0x61, 0xbc,
	0xa8, 0xf5, 0x08, 0x4e, 0x9d, 0xd1, 0xf1, 0xb2, 0x10, 0x0e, 0xd4, 0x40, 0xc8, 0x6c, 0xa0, 0x15,
	0xee, 0x42, 0xdc, 0x42, 0x86, 0x92, 0xc5, 0x4d, 0x68, 0xed, 0x67, 0xf1, 0xc6, 0x2f, 0xe1, 0x27,
	0xf2, 0x8b, 0x38, 0xef, 0x6e, 0x28, 0xd0, 0x76, 0xe6, 0x5c, 0xf1, 0x7b, 0x9e, 0xcd, 0xbe, 0xef,
	0x9b, 0x87, 0xec, 0x92, 0x93, 0xb5, 0x14, 0xb1, 0x88, 0x9a, 0x0f, 0x5e, 0xec, 0x2f, 0xb8, 0x6c,
	0x28, 0x59, 0xbf, 0x26, 0xe5, 0x51, 0xe0, 0x2f, 0xb9, 0x64, 0xfc, 0xcf, 0x0d, 0x8f, 0x62, 0xfa,
	0x96, 0xe4, 0xb4, 0x61, 0xa4, 0x6a, 0xa9, 0xb3, 0x22, 0x4b, 0x14, 0xfd, 0x81, 0x94, 0x4c, 0x1e,
	0xc5, 0x41, 0xe8, 0xc5, 0x81, 0x08, 0x8d, 0x74, 0x2d, 0x75, 0x56, 0xb9, 0x28, 0x35, 0xda, 0x1b,
	0x29, 0x79, 0xe8, 0x07, 0x3c, 0x62, 0xfb, 0xeb, 0xf5, 0x7f, 0xd2, 0xa4, 0xb2, 0x2d, 0x1c, 0xad,
	0x45, 0x18, 0x71, 0xac, 0xec, 0x3e, 0xae, 0x66, 0xe2, 0x6e, 0x5b, 0x59, 0x2b, 0x4a, 0x49, 0x76,
	0xb8, 0xe6, 0xba, 0x64, 0x91, 0x29, 0x46, 0xaf, 0x17, 0xcc, 0x17, 0x46, 0x46, 0x7b, 0xc8, 0x14,
	0x48, 0xa6, 0x2f, 0x1e, 0x8c, 0xac, 0xb2, 0x10, 0xe9, 0x09, 0x39, 0x72, 0x64, 0xe0, 0x73, 0xe3,
	0x48, 0x79, 0x5a, 0xd0, 0xda, 0xe1, 0xa4, 0x39, 0xb5, 0xb6, 0x6f, 0xe1, 0x24, 0xd7, 0xe2, 0x6e,
	0xb3, 0xe2, 0x46, 0x5e, 0x4f, 0xa2, 0x95, 0xea, 0x30, 0x32, 0x8d, 0x42, 0xd2, 0x61, 0x64, 0xd2,
	0x77, 0xa4, 0xe8, 0x48, 0x7e, 0xdf, 0xbe, 0x13, 0x11, 0x37, 0x8a, 0xca, 0xdf, 0x19, 0x58, 0xa7,
	0xbd, 0xf0, 0xc2, 0x39, 0x37, 0x88, 0xae, 0xa3, 0x15, 0xfd, 0x96, 0x94, 0x1d, 0x2e, 0x7d, 0x1e,
	0xc6, 0xc9, 0x72, 0x49, 0x2d, 0x1f, 0x9a, 0xf5, 0xff, 0x32, 0xa4, 0xda, 0x16, 0xab, 0xb5, 0x17,
	0x3e, 0xee, 0x67, 0xf4, 0x6a, 0xfa, 0x94, 0x64, 0x6d, 0x6f, 0xc5, 0xb7, 0x19, 0x21, 0xd3, 0x53,
	0x52, 0xe8, 0xfc, 0xe5, 0xeb, 0x06, 0x3a, 0xa7, 0x27, 0xad, 0xb2, 0xe6, 0x7e, 0x2c, 0x64, 0x12,
	0x57, 0xa2, 0xf0, 0x7d, 0x06, 0x9e, 0x5c, 0xf2, 0xb8, 0xed, 0xad, 0x93, 0xd4, 0x76, 0x06, 0x35,
	0x48, 0xde, 0xe9, 0x30, 0xcc, 0x28, 0x49, 0x6d, 0x2b, 0xb1, 0x97, 0xd3, 0xe9, 0xea, 0x25, 0x9d,
	0xd9, 0x93, 0x56, 0x79, 0x07, 0xf7, 0x0e, 0x97, 0xee, 0xc2, 0x93, 0x3c, 0x49, 0x6f, 0xdf, 0xc2,
	0x5c, 0x3b, 0x8e, 0x9b, 0xe4, 0x87, 0x88, 0x7b, 0x18, 0xdf, 0xed, 0xd1, 0xf1, 0xed, 0x5b, 0xb4,
	0x4e, 0x8e, 0x1d, 0x29, 0x6e, 0x83, 0x78, 0xe0, 0xc9, 0x79, 0x10, 0x26, 0x11, 0x1e, 0x78, 0x38,
	0xd5, 0x84, 0x7b, 0x52, 0x7d, 0x29, 0xc7, 0x7a, 0xaa, 0xad, 0xc6, 0x77, 0x41, 0xc6, 0x2f, 0xa6,
	0xac, 0xdf, 0x25, 0x91, 0xf4, 0x3d, 0x79, 0xa3, 0x5a, 0x44, 0xc3, 0x4d, 0x1c, 0xc5, 0x5e, 0xf8,
	0x47, 0x10, 0xce, 0x8d, 0x8a, 0x7a, 0xe6, 0xe5, 0x02, 0x3d, 0x27, 0xa0, 0x3e, 0xab, 0x91, 0x68,
	0x09, 0xb1, 0xd4, 0x09, 0x54, 0xd5, 0xc3, 0x2f, 0x7c, 0xfc, 0x97, 0x5a, 0x3c, 0xf6, 0x0c, 0xd0,
	0xff, 0x12, 0xf2, 0xf9, 0xbf, 0x69, 0x42, 0x76, 0x87, 0x84, 0xe6, 0x49, 0x66, 0xec, 0x9a, 0xf0,
	0x09, 0x42, 0x67, 0xcc, 0x20, 0x85, 0xf0, 0x9b, 0x33, 0x81, 0x34, 0x42, 0xab, 0x6b, 0x43, 0x06,
	0xa1, 0x3d, 0xb5, 0x20, 0x8b, 0x60, 0x5a, 0x16, 0x1c, 0x21, 0x74, 0x5b, 0x0e, 0xe4, 0x10, 0x7a,
	0xe3, 0x0f, 0x90, 0x47, 0x70, 0xfa, 0x36, 0x14, 0x10, 0xd8, 0xd0, 0x86, 0x22, 0x82, 0xdb, 0xb1,
	0x80, 0xa8, 0xed, 0xbd, 0x0f, 0x50, 0x42, 0xb8, 0x72, 0x2d, 0x38, 0x46, 0xb0, 0x87, 0x16, 0x94,
	0xd5, 0x76, 0x66, 0x41, 0x45, 0xed, 0x1a, 0xb7, 0xa0, 0x8a, 0x30, 0x62, 0x13, 0x00, 0x84, 0xcb,
	0xb1, 0x09, 0x6f, 0xd4, 0x18, 0xac, 0x0f, 0x54, 0xd5, 0xb9, 0x34, 0xe1, 0x53, 0x05, 0xf6, 0x04,
	0x4e, 0xd4, 0x76, 0xcb, 0x84, 0xcf, 0x54, 0x65, 0x93, 0xc1, 0x5b, 0x05, 0x7d, 0x17, 0x3e, 0x57,
	0x60, 0x33, 0x30, 0x10, 0x2c, 0x76, 0x03, 0x5f, 0x20, 0x0c, 0x7e, 0xb7, 0xe1, 0x54, 0xc1, 0x84,
	0xc1, 0x97, 0x6a, 0x8c, 0xa9, 0x09, 0xef, 0xd4, 0xf0, 0x3d, 0x07, 0xbe, 0x52, 0x33, 0x77, 0x4d,
	0xf8, 0x5a, 0x8d, 0xd1, 0x6b, 0xc1, 0x37, 0x08, 0xd3, 0x4b, 0x06, 0xb5, 0x8b, 0xbf, 0x53, 0x24,
	0x7f, 0xa3, 0xaf, 0x2a, 0x7a, 0x4e, 0xf2, 0x5d, 0x1e, 0x5f, 0x85, 0xb7, 0x82, 0x56, 0x1a, 0x07,
	0xd7, 0xd5, 0x69, 0xb5, 0xf1, 0xec, 0x96, 0x79, 0x4f, 0x0a, 0x03, 0x21, 0xf9, 0xab, 0x0f, 0x43,
	0xe3, 0xf9, 0x79, 0xfb, 0x99, 0x54, 0xdd, 0xcd, 0x2c, 0xf2, 0x65, 0x30, 0xe3, 0xc9, 0x51, 0xfb,
	0x58, 0x87, 0xb3, 0xd4, 0x8f, 0xa9, 0xd6, 0xf7, 0xd3, 0xef, 0xe6, 0x41, 0xbc, 0xd8, 0xcc, 0x1a,
	0xbe, 0x58, 0x35, 0x6f, 0xf1, 0x5c, 0xdf, 0x07, 0xcb, 0x5f, 0x2e, 0x7e, 0x6d, 0xba, 0xa3, 0xbe,
	0xd5, 0x61, 0xcd, 0xb9, 0x5c, 0xfb, 0x4d, 0x7d, 0xe3, 0xce, 0x72, 0xea, 0xf7, 0xa7, 0xff, 0x03,
	0x00, 0x00, 0xff, 0xff, 0xe3, 0x50, 0x94, 0x05, 0x82, 0x05, 0x00, 0x00,
}
