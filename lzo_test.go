package lzo

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

type lzoTest struct {
	name string
	desc string
	raw  string
	lzo  []byte
	err  error
}

var lzoTests = []lzoTest{
	{
		"empty.txt",
		"empty.txt",
		"",
		[]byte{
			0x89, 0x4c, 0x5a, 0x4f, 0x0, 0xd, 0xa, 0x1a, 0xa,
			0x10, 0x30, 0x20, 0x60, 0x9, 0x40, 0x1, 0x5, 0x3,
			0x0, 0x0, 0x1, 0x0, 0x0, 0x81, 0xb4, 0x51, 0xcf,
			0x50, 0x65, 0x0, 0x0, 0x0, 0x0, 0x9, 0x65, 0x6d,
			0x70, 0x74, 0x79, 0x2e, 0x74, 0x78, 0x74, 0x6a,
			0x20, 0x7, 0xe4, 0x0, 0x0, 0x0, 0x0,
		},
		nil,
	},
	{
		"hello.txt",
		"hello.txt",
		"hello world\n",
		[]byte{
			0x89, 0x4c, 0x5a, 0x4f, 0x0, 0xd, 0xa, 0x1a, 0xa,
			0x10, 0x30, 0x20, 0x60, 0x9, 0x40, 0x1, 0x5, 0x3,
			0x0, 0x0, 0x1, 0x0, 0x0, 0x81, 0xa4, 0x51, 0xcf,
			0x8a, 0xf, 0x0, 0x0, 0x0, 0x0, 0x9, 0x68, 0x65,
			0x6c, 0x6c, 0x6f, 0x2e, 0x74, 0x78, 0x74, 0x66,
			0xe3, 0x7, 0x9d, 0x0, 0x0, 0x0, 0xc, 0x0, 0x0,
			0x0, 0xc, 0x1e, 0x72, 0x4, 0x67, 0x68, 0x65,
			0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c,
			0x64, 0xa, 0x0, 0x0, 0x0, 0x0,
		},
		nil,
	},
	{
		"hello.txt",
		"hello.txt x2",
		"hello world\n" +
			"hello world\n",
		[]byte{
			0x89, 0x4c, 0x5a, 0x4f, 0x0, 0xd, 0xa, 0x1a, 0xa,
			0x10, 0x30, 0x20, 0x60, 0x9, 0x40, 0x1, 0x5, 0x3,
			0x0, 0x0, 0x1, 0x0, 0x0, 0x81, 0xa4, 0x51, 0xd6,
			0x36, 0xf8, 0x0, 0x0, 0x0, 0x0, 0x9, 0x68, 0x65,
			0x6c, 0x6c, 0x6f, 0x2e, 0x74, 0x78, 0x74, 0x6f,
			0xc1, 0x8, 0x39, 0x0, 0x0, 0x0, 0x18, 0x0, 0x0,
			0x0, 0x18, 0x71, 0xac, 0x8, 0xcd, 0x68, 0x65, 0x6c,
			0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
			0xa, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77,
			0x6f, 0x72, 0x6c, 0x64, 0xa, 0x0, 0x0, 0x0, 0x0,
		},
		nil,
	},
	{
		"shesells.txt",
		"shesells.txt",
		"she sells seashells by the seashore\n",
		[]byte{
			0x89, 0x4c, 0x5a, 0x4f, 0x0, 0xd, 0xa, 0x1a, 0xa,
			0x10, 0x30, 0x20, 0x60, 0x9, 0x40, 0x1, 0x5, 0x3,
			0x0, 0x0, 0x1, 0x0, 0x0, 0x81, 0xb6, 0x51, 0xd6,
			0x38, 0x38, 0x0, 0x0, 0x0, 0x0, 0xc, 0x73, 0x68,
			0x65, 0x73, 0x65, 0x6c, 0x6c, 0x73, 0x2e, 0x74,
			0x78, 0x74, 0x80, 0x24, 0x8, 0xdf, 0x0, 0x0, 0x0,
			0x24, 0x0, 0x0, 0x0, 0x24, 0xfa, 0x2c, 0xd, 0x48,
			0x73, 0x68, 0x65, 0x20, 0x73, 0x65, 0x6c, 0x6c,
			0x73, 0x20, 0x73, 0x65, 0x61, 0x73, 0x68, 0x65,
			0x6c, 0x6c, 0x73, 0x20, 0x62, 0x79, 0x20, 0x74,
			0x68, 0x65, 0x20, 0x73, 0x65, 0x61, 0x73, 0x68,
			0x6f, 0x72, 0x65, 0xa, 0x0, 0x0, 0x0, 0x0,
		},
		nil,
	},
	{
		"gettysburg.txt",
		"gettysburg",
		"  Four score and seven years ago our fathers brought forth on\n" +
			"this continent, a new nation, conceived in Liberty, and dedicated\n" +
			"to the proposition that all men are created equal.\n" +
			"  Now we are engaged in a great Civil War, testing whether that\n" +
			"nation, or any nation so conceived and so dedicated, can long\n" +
			"endure.\n" +
			"  We are met on a great battle-field of that war.\n" +
			"  We have come to dedicate a portion of that field, as a final\n" +
			"resting place for those who here gave their lives that that\n" +
			"nation might live.  It is altogether fitting and proper that\n" +
			"we should do this.\n" +
			"  But, in a larger sense, we can not dedicate — we can not\n" +
			"consecrate — we can not hallow — this ground.\n" +
			"  The brave men, living and dead, who struggled here, have\n" +
			"consecrated it, far above our poor power to add or detract.\n" +
			"The world will little note, nor long remember what we say here,\n" +
			"but it can never forget what they did here.\n" +
			"  It is for us the living, rather, to be dedicated here to the\n" +
			"unfinished work which they who fought here have thus far so\n" +
			"nobly advanced.  It is rather for us to be here dedicated to\n" +
			"the great task remaining before us — that from these honored\n" +
			"dead we take increased devotion to that cause for which they\n" +
			"gave the last full measure of devotion —\n" +
			"  that we here highly resolve that these dead shall not have\n" +
			"died in vain — that this nation, under God, shall have a new\n" +
			"birth of freedom — and that government of the people, by the\n" +
			"people, for the people, shall not perish from this earth.\n" +
			"\n" +
			"Abraham Lincoln, November 19, 1863, Gettysburg, Pennsylvania\n",
		[]byte{
			0x89, 0x4c, 0x5a, 0x4f, 0x0, 0xd, 0xa, 0x1a, 0xa, 0x10, 0x30, 0x20, 0x60,
			0x9, 0x40, 0x1, 0x5, 0x3, 0x0, 0x0, 0x1, 0x0, 0x0, 0x81, 0xb6, 0x51, 0xd7,
			0x65, 0x87, 0x0, 0x0, 0x0, 0x0, 0xe, 0x67, 0x65, 0x74, 0x74, 0x79, 0x73,
			0x62, 0x75, 0x72, 0x67, 0x2e, 0x74, 0x78, 0x74, 0x9e, 0x3b, 0xa, 0x4b,
			0x0, 0x0, 0x6, 0x1a, 0x0, 0x0, 0x4, 0xa6, 0x18, 0x31, 0x2b, 0x27, 0x0,
			0x5f, 0x20, 0x20, 0x46, 0x6f, 0x75, 0x72, 0x20, 0x73, 0x63, 0x6f, 0x72,
			0x65, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x73, 0x65, 0x76, 0x65, 0x6e, 0x20,
			0x79, 0x65, 0x61, 0x72, 0x73, 0x20, 0x61, 0x67, 0x6f, 0x20, 0x6f, 0x75,
			0x72, 0x20, 0x66, 0x61, 0x74, 0x68, 0x65, 0x72, 0x73, 0x20, 0x62, 0x72,
			0x6f, 0x75, 0x67, 0x68, 0x74, 0x20, 0x66, 0x6f, 0x72, 0x74, 0x68, 0x20,
			0x6f, 0x6e, 0xa, 0x74, 0x68, 0x69, 0x73, 0x20, 0x63, 0x6f, 0x6e, 0x74,
			0x69, 0x6e, 0x65, 0x6e, 0x74, 0x2c, 0x20, 0x61, 0x20, 0x6e, 0x65, 0x77,
			0x20, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2c, 0x20, 0x63, 0x6f, 0x6e,
			0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x20, 0x69, 0x6e, 0x20, 0x4c, 0x69,
			0x62, 0x65, 0x72, 0x74, 0x79, 0x2c, 0x90, 0xc, 0x0, 0x34, 0x64, 0x65,
			0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0xa, 0x74, 0x6f, 0x20, 0x74,
			0x68, 0x65, 0x20, 0x70, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69,
			0x6f, 0x6e, 0x20, 0x74, 0x68, 0x61, 0x74, 0x20, 0x61, 0x6c, 0x6c, 0x20,
			0x6d, 0x65, 0x6e, 0x20, 0x61, 0x72, 0x65, 0x20, 0x63, 0x72, 0x65, 0x61,
			0x74, 0x65, 0x64, 0x20, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x2e, 0xa, 0x20,
			0x20, 0x4e, 0x6f, 0x77, 0x20, 0x77, 0x65, 0x20, 0x6c, 0x3, 0x4, 0x65,
			0x6e, 0x67, 0x61, 0x67, 0x65, 0x64, 0x64, 0xc, 0x0, 0x10, 0x61, 0x20,
			0x67, 0x72, 0x65, 0x61, 0x74, 0x20, 0x43, 0x69, 0x76, 0x69, 0x6c, 0x20,
			0x57, 0x61, 0x72, 0x2c, 0x20, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67,
			0x20, 0x77, 0x68, 0x65, 0x74, 0x68, 0x65, 0x72, 0x8a, 0xb, 0xa, 0x6e,
			0xd8, 0x13, 0x4, 0x6f, 0x72, 0x20, 0x61, 0x6e, 0x79, 0x20, 0xbb, 0x1,
			0x20, 0x73, 0x6f, 0x29, 0xbc, 0x2, 0x82, 0x21, 0x6f, 0x20, 0x27, 0x98,
			0x2, 0x0, 0x6, 0x2c, 0x20, 0x63, 0x61, 0x6e, 0x20, 0x6c, 0x6f, 0x6e,
			0x67, 0xa, 0x65, 0x6e, 0x64, 0x75, 0x72, 0x65, 0x2e, 0xa, 0x20, 0x20,
			0x57, 0x65, 0x20, 0x64, 0x10, 0x4, 0x6d, 0x65, 0x74, 0x20, 0x6f, 0x6e,
			0x20, 0xf4, 0xf, 0xc, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x2d, 0x66,
			0x69, 0x65, 0x6c, 0x64, 0x20, 0x6f, 0x66, 0x88, 0xe, 0x1, 0x20, 0x77,
			0x61, 0x72, 0xc4, 0x6, 0x8, 0x68, 0x61, 0x76, 0x65, 0x20, 0x63, 0x6f,
			0x6d, 0x65, 0x20, 0x74, 0x28, 0x7c, 0x1, 0x3, 0x20, 0x61, 0x20, 0x70,
			0x6f, 0x72, 0x90, 0x1f, 0xe4, 0x6, 0x9c, 0x7, 0xb, 0x2c, 0x20, 0x61,
			0x73, 0x20, 0x61, 0x20, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0xa, 0x72, 0xcc,
			0x19, 0x2, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x68, 0x30, 0xe, 0x20, 0x74,
			0x68, 0x6f, 0x73, 0x65, 0x20, 0x77, 0x68, 0x6f, 0x20, 0x68, 0x65, 0x72,
			0x65, 0x20, 0x67, 0x68, 0xb, 0x8, 0x74, 0x68, 0x65, 0x69, 0x72, 0x20,
			0x6c, 0x69, 0x76, 0x65, 0x73, 0xa8, 0xf, 0x74, 0x29, 0xc8, 0x1e, 0x3,
			0x20, 0x6d, 0x69, 0x67, 0x68, 0x74, 0x90, 0x3, 0xb, 0x2e, 0x20, 0x20,
			0x49, 0x74, 0x20, 0x69, 0x73, 0x20, 0x61, 0x6c, 0x74, 0x6f, 0x67, 0xb3,
			0x23, 0x66, 0x69, 0x74, 0x90, 0x25, 0x70, 0x1f, 0x7c, 0x31, 0xfc, 0x25,
			0x8, 0x77, 0x65, 0x20, 0x73, 0x68, 0x6f, 0x75, 0x6c, 0x64, 0x20, 0x64,
			0x72, 0x35, 0x69, 0x73, 0x68, 0x19, 0x1, 0x42, 0x75, 0x74, 0x2c, 0xb0,
			0x2e, 0xd, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x72, 0x20, 0x73, 0x65, 0x6e,
			0x73, 0x65, 0x2c, 0x20, 0x77, 0x65, 0x8b, 0x25, 0x6e, 0x6f, 0x74, 0x27,
			0xf4, 0x4, 0x1, 0x20, 0xe2, 0x80, 0x94, 0x29, 0x5c, 0x0, 0x8, 0xa, 0x63,
			0x6f, 0x6e, 0x73, 0x65, 0x63, 0x72, 0x61, 0x74, 0x65, 0x2d, 0x64, 0x0,
			0x4, 0x20, 0x68, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x94, 0x2, 0x84, 0x4b,
			0x3, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x78, 0xd, 0x3, 0x54, 0x68,
			0x65, 0x20, 0x62, 0x72, 0x60, 0x1c, 0x1, 0x6d, 0x65, 0x6e, 0x2c, 0x68,
			0x18, 0x78, 0x3a, 0x64, 0x15, 0x2, 0x64, 0x65, 0x61, 0x64, 0x2c, 0x90,
			0x20, 0x6, 0x73, 0x74, 0x72, 0x75, 0x67, 0x67, 0x6c, 0x65, 0x64, 0x9a,
			0x21, 0x2c, 0x20, 0x68, 0x2d, 0x29, 0xb0, 0x1, 0xc, 0x64, 0x20, 0x69,
			0x74, 0x2c, 0x20, 0x66, 0x61, 0x72, 0x20, 0x61, 0x62, 0x6f, 0x76, 0x65,
			0x80, 0x5b, 0x5, 0x70, 0x6f, 0x6f, 0x72, 0x20, 0x70, 0x6f, 0x77, 0x68,
			0x1d, 0xf, 0x6f, 0x20, 0x61, 0x64, 0x64, 0x20, 0x6f, 0x72, 0x20, 0x64,
			0x65, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e, 0xa, 0x70, 0xe, 0xa, 0x77,
			0x6f, 0x72, 0x6c, 0x64, 0x20, 0x77, 0x69, 0x6c, 0x6c, 0x20, 0x6c, 0x69,
			0x64, 0x3b, 0x7, 0x20, 0x6e, 0x6f, 0x74, 0x65, 0x2c, 0x20, 0x6e, 0x6f,
			0x72, 0x9c, 0x41, 0xb, 0x20, 0x72, 0x65, 0x6d, 0x65, 0x6d, 0x62, 0x65,
			0x72, 0x20, 0x77, 0x68, 0x61, 0x74, 0x7b, 0x1c, 0x73, 0x61, 0x79, 0xa0,
			0x10, 0x4, 0xa, 0x62, 0x75, 0x74, 0x20, 0x69, 0x74, 0xbc, 0x21, 0x1,
			0x65, 0x76, 0x65, 0x72, 0x73, 0x36, 0x67, 0x65, 0x74, 0xa4, 0x5, 0x4,
			0x74, 0x68, 0x65, 0x79, 0x20, 0x64, 0x69, 0xb0, 0x15, 0x68, 0x1c, 0xbc,
			0x31, 0x2, 0x66, 0x6f, 0x72, 0x20, 0x75, 0x69, 0x37, 0x65, 0xd3, 0x1c,
			0x2c, 0x20, 0x72, 0x95, 0x6f, 0x2c, 0x62, 0x46, 0x62, 0x65, 0x27, 0x54,
			0x5, 0xa0, 0x7, 0x70, 0x2, 0x0, 0x7, 0x74, 0x68, 0x65, 0xa, 0x75, 0x6e,
			0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x20, 0x77, 0x6f, 0x72,
			0x6b, 0x20, 0x77, 0x68, 0x69, 0x63, 0x68, 0x65, 0x6b, 0x79, 0x8a, 0x23,
			0x66, 0x6f, 0x9c, 0x76, 0x98, 0x44, 0x78, 0x22, 0x2, 0x20, 0x74, 0x68,
			0x75, 0x73, 0x8c, 0x21, 0xe, 0x73, 0x6f, 0xa, 0x6e, 0x6f, 0x62, 0x6c,
			0x79, 0x20, 0x61, 0x64, 0x76, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x27, 0x64,
			0x8, 0xb8, 0xe, 0x6c, 0x16, 0x80, 0x12, 0x1, 0x6f, 0x20, 0x62, 0x65,
			0x9c, 0x1a, 0x29, 0x7, 0x2, 0x74, 0x6f, 0xa, 0x70, 0x77, 0xb8, 0x6e,
			0x1, 0x74, 0x61, 0x73, 0x6b, 0x77, 0x21, 0x61, 0x69, 0x6e, 0x6f, 0x33,
			0x62, 0x65, 0x66, 0x6a, 0x8a, 0x75, 0x73, 0xd8, 0x39, 0x4, 0x61, 0x74,
			0x20, 0x66, 0x72, 0x6f, 0x6d, 0x6c, 0x12, 0x8, 0x73, 0x65, 0x20, 0x68,
			0x6f, 0x6e, 0x6f, 0x72, 0x65, 0x64, 0xa, 0x74, 0x37, 0x70, 0x26, 0x4,
			0x74, 0x61, 0x6b, 0x65, 0x20, 0x69, 0x6e, 0x64, 0x7d, 0x5, 0x73, 0x65,
			0x64, 0x20, 0x64, 0x65, 0x76, 0x6f, 0x98, 0x61, 0x99, 0x83, 0x61, 0x6e,
			0x28, 0x75, 0x73, 0xb0, 0x5e, 0x28, 0x65, 0x3, 0xa, 0xe0, 0x5e, 0x6,
			0x20, 0x6c, 0x61, 0x73, 0x74, 0x20, 0x66, 0x75, 0x6c, 0x68, 0x86, 0x2,
			0x61, 0x73, 0x75, 0x72, 0x65, 0x60, 0x6f, 0x27, 0xc, 0x1, 0x2, 0xe2,
			0x80, 0x94, 0xa, 0x20, 0xb1, 0x61, 0x77, 0xc8, 0x18, 0xb, 0x68, 0x69,
			0x67, 0x68, 0x6c, 0x79, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65,
			0xac, 0x3, 0x1, 0x74, 0x68, 0x65, 0x73, 0x78, 0x2b, 0x1, 0x61, 0x64,
			0x20, 0x73, 0x60, 0x51, 0x60, 0x3d, 0xae, 0x48, 0x64, 0x69, 0xa0, 0x8c,
			0x1, 0x76, 0x61, 0x69, 0x6e, 0x28, 0x34, 0x3, 0x84, 0x54, 0xb0, 0x87,
			0x9, 0x2c, 0x20, 0x75, 0x6e, 0x64, 0x65, 0x72, 0x20, 0x47, 0x6f, 0x64,
			0x2c, 0xd0, 0x7, 0x95, 0x2c, 0x61, 0x78, 0xa1, 0x1, 0xa, 0x62, 0x69,
			0x72, 0x74, 0xa5, 0x6, 0x66, 0x20, 0x66, 0x72, 0x65, 0x65, 0x64, 0x6f,
			0x6d, 0x88, 0x8, 0x68, 0x57, 0x68, 0x73, 0x8, 0x20, 0x67, 0x6f, 0x76,
			0x65, 0x72, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x16, 0x78, 0x29, 0x7,
			0x70, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x2c, 0x20, 0x62, 0x79, 0x71, 0x25,
			0xa, 0xf8, 0x1, 0x64, 0x43, 0x2a, 0x78, 0x0, 0x28, 0x90, 0x2, 0x3, 0x70,
			0x65, 0x72, 0x69, 0x73, 0x68, 0xf8, 0x2b, 0x0, 0xd, 0x69, 0x73, 0x20,
			0x65, 0x61, 0x72, 0x74, 0x68, 0x2e, 0xa, 0xa, 0x41, 0x62, 0x72, 0x61,
			0x68, 0x61, 0x6d, 0x20, 0x4c, 0x69, 0x6e, 0x63, 0x6f, 0x6c, 0x6e, 0x2c,
			0x20, 0x4e, 0x6f, 0x76, 0xb0, 0x55, 0x0, 0x11, 0x31, 0x39, 0x2c, 0x20,
			0x31, 0x38, 0x36, 0x33, 0x2c, 0x20, 0x47, 0x65, 0x74, 0x74, 0x79, 0x73,
			0x62, 0x75, 0x72, 0x67, 0x2c, 0x20, 0x50, 0x65, 0x6e, 0x6e, 0x73, 0x79,
			0x6c, 0x76, 0x61, 0x6e, 0x69, 0x61, 0xa, 0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		},
		nil,
	},
}

func TestCompressor(t *testing.T) {
	for _, tt := range lzoTests {
		in := strings.NewReader(tt.raw)
		buf := new(bytes.Buffer)
		lzo := NewWriter(buf)
		_, err := io.Copy(lzo, in)
		if err != nil {
			t.Errorf("%s: Write: %s", tt.name, err)
		}
		if !bytes.Equal(buf.Bytes(), tt.lzo) {
			t.Errorf("%s: got %#v want %#v", tt.name, buf.Bytes(), tt.lzo)
		}
	}
}

func TestDecompressor(t *testing.T) {
	b := new(bytes.Buffer)
	for _, tt := range lzoTests {
		in := bytes.NewBuffer(tt.lzo)
		lzo, err := NewReader(in)
		if err != nil {
			t.Errorf("%s: NewReader: %s", tt.name, err)
			continue
		}
		defer lzo.Close()
		if tt.name != lzo.Name {
			t.Errorf("%s: got name %s", tt.name, lzo.Name)
		}
		b.Reset()
		n, err := io.Copy(b, lzo)
		if err != tt.err {
			t.Errorf("%s: io.Copy: %v want %v", tt.name, err, tt.err)
		}
		s := b.String()
		if s != tt.raw {
			t.Errorf("%s: got %d-byte %q want %d-byte %q", tt.name, n, s, len(tt.raw), tt.raw)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	w.Name = "name"
	if _, err := w.Write([]byte("payload")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	r, err := NewReader(buf)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(b) != "payload" {
		t.Fatalf("payload is %q, want %q", string(b), "payload")
	}
	if r.Name != "name" {
		t.Fatalf("name is %q, want %q", r.Name, "name")
	}
	if err := r.Close(); err != nil {
		t.Fatalf("Reader.Close: %v", err)
	}
}
