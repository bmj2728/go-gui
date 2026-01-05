package api

import (
	"net/url"
)

type CAASImageType int

const (
	CAASImageTypeSquare CAASImageType = iota
	CAASImageTypeMedium
	CAASImageTypeSmall
	CAASImageTypeXSmall
)

var CAASImageTypes = map[CAASImageType]string{
	CAASImageTypeSquare: "square",
	CAASImageTypeMedium: "medium",
	CAASImageTypeSmall:  "small",
	CAASImageTypeXSmall: "xsmall",
}

type CAASImageFilter int

const (
	CAASImageFilterMono CAASImageFilter = iota
	CAASImageFilterNegate
	CAASImageFilterCustom
)

var CAASImageFilters = map[CAASImageFilter]string{
	CAASImageFilterMono:   "mono",
	CAASImageFilterNegate: "negate",
	CAASImageFilterCustom: "custom",
}

type CAASImageFit int

const (
	CAASImageFitCover CAASImageFit = iota
	CAASImageFitContain
	CAASImageFitFill
	CAASImageFitInside
	CAASImageFitOutside
)

var CAASImageFits = map[CAASImageFit]string{
	CAASImageFitCover:   "cover",
	CAASImageFitContain: "contain",
	CAASImageFitFill:    "fill",
	CAASImageFitInside:  "inside",
	CAASImageFitOutside: "outside",
}

type CAASImagePosition int

const (
	CAASImagePositionCenter CAASImagePosition = iota
	CAASImagePositionTop
	CAASImagePositionRightTop
	CAASImagePositionRight
	CAASImagePositionRightBottom
	CAASImagePositionBottom
	CAASImagePositionLeftBottom
	CAASImagePositionLeft
	CAASImagePositionLeftTop
)

var CAASImagePositions = map[CAASImagePosition]string{
	CAASImagePositionCenter:      "center",
	CAASImagePositionTop:         "top",
	CAASImagePositionRightTop:    url.QueryEscape("right top"), // right%20top
	CAASImagePositionRight:       "right",
	CAASImagePositionRightBottom: url.QueryEscape("right bottom"), // right%20bottom
	CAASImagePositionBottom:      "bottom",
	CAASImagePositionLeftBottom:  url.QueryEscape("left bottom"), // left%20top
	CAASImagePositionLeft:        "left",
	CAASImagePositionLeftTop:     url.QueryEscape("left top"), // left%20top
}

type CAASFont int

const (
	CAASFontImpact CAASFont = iota
	CAASFontAndale
	CAASFontMono
	CAASFontArial
	CAASFontArialBlack
	CAASFontComicSansMS
	CAASFontCourierNew
	CAASFontGeorgia
	CAASFontTimesNewRoman
	CAASFontVerdana
	CAASFontWebdings
)

var CAASFonts = map[CAASFont]string{
	CAASFontImpact:        "Impact",
	CAASFontAndale:        "Andale",
	CAASFontMono:          "Mono",
	CAASFontArial:         "Arial",
	CAASFontArialBlack:    url.QueryEscape("Arial Black"),
	CAASFontComicSansMS:   url.QueryEscape("Comic Sans MS"),
	CAASFontCourierNew:    url.QueryEscape("Courier New"),
	CAASFontGeorgia:       "Georgia",
	CAASFontTimesNewRoman: url.QueryEscape("Times New Roman"),
	CAASFontVerdana:       "Verdana",
	CAASFontWebdings:      "Webdings",
}
