package gosketch

import (
	"errors"
	"strconv"
)

type Css struct {
	Value []PageCss
}

type PageCss struct {
	ID     string
	Struct []BlockCss
}

type BlockCss struct {
	Width           float64
	Height          float64
	Top             float64
	Left            float64
	BackgroundColor string
	BackgroundImage string
	BorderRadius    float64
	Borders         []string
	BoxShadow       string
	Children        []BlockCss
	Font            interface{}
}

type Font struct {
	Size    float64
	Color   string
	Family  string
	Weight  float64
	Spacing float64
	Height  float64
	Text    string
}

type MapShadow struct {
	Value map[string]interface{}
}

func (s *SketchFile) GetCSS() *Css {
	result := make([]PageCss, 0)
	for key, page := range s.Pages {
		newPage := PageCss{ID: key, Struct: make([]BlockCss, len(page.Layers))}
		countW := len(page.Layers)
		countWoods := make(chan int)
		growBrancge := make(chan int)
		count := len(page.Layers)
		for index, item := range page.Layers {
			go cssBlock(item, index, &newPage.Struct[index], countWoods, growBrancge)
		}
		for countW > 0 {
			select {
			case c := <-countWoods:
				countW = countW + c - 1
			case s := <-growBrancge:
				count = count - s
			}
		}
		close(countWoods)
		close(growBrancge)
		result = append(result, newPage)
	}
	return &Css{Value: result}
}

func cssBlock(layer map[string]interface{}, index int, block *BlockCss, countWoods chan<- int, growBranche chan<- int) {
	frame, okF := layer["frame"].(map[string]interface{})
	if okF {
		block.getPosition(frame)
	}
	bkgM, ok := layer["backgroundColor"].(map[string]interface{})
	if ok {
		block.BackgroundColor = colorRGBA(bkgM)
	}
	style, ok := layer["style"].(map[string]interface{})
	if ok {
		shadowS, ok := style["shadow"].([]map[string]interface{})
		if ok {
			for index, item := range shadowS {
				shd := MapShadow{Value: item}
				itemShadow, err := shd.boxShadow()
				if err == nil {
					if index > 0 {
						block.BoxShadow = block.BoxShadow + ", "
					}
					block.BoxShadow = block.BoxShadow + itemShadow
				}
			}
		}
		borders, ok := style["borders"].([]interface{})
		if ok {
			go block.getBorders(borders)
		}
	}
	if layer["_class"] == "artboard" || layer["_class"] == "group" || layer["_class"] == "shapeGroup" || layer["_class"] == "symbolMaster" {
		block.Font = nil
	} else {
		atributeString, ok := layer["attributedString"].(map[string]interface{})
		if ok {
			_, ok := atributeString["archivedAttributedString"].(map[string]interface{})
			if !ok {
				block.fontStyle(atributeString)
			}
		}
	}
	childrenMaps, ok := layer["layers"].([]interface{})
	growBranche <- 1
	if ok {
		block.getChildren(childrenMaps, countWoods)
	} else {
		countWoods <- 0
	}
}

func (block *BlockCss) getChildren(childrenMaps []interface{}, countWoods chan<- int) {
	structureBranches := make([]BlockCss, len(childrenMaps))
	growBranch := make(chan int)
	count := len(childrenMaps)
	countWoods <- len(childrenMaps)
	for index, child := range childrenMaps {
		child, ok := child.(map[string]interface{})
		if ok {
			go cssBlock(child, index, &structureBranches[index], countWoods, growBranch)
		}
	}
	for count > 0 {
		select {
		case s := <-growBranch:
			count = count - s
		}
	}
	close(growBranch)
	block.Children = structureBranches
}

func (block *BlockCss) getBorders(borders []interface{}) {
	result := make([]string, 0)
	for _, border := range borders {
		border, ok := border.(map[string]interface{})
		if ok && border["isEnabled"].(bool) {
			width := strconv.Itoa(int(border["thickness"].(float64)))
			color, ok := border["color"].(map[string]interface{})
			if ok {
				colorString := colorRGBA(color)
				result = append(result, width+"px solid "+colorString) // TODO: only solid
			}
		}
	}
	if len(result) > 0 {
		block.Borders = result
	}
}

func (block *BlockCss) getPosition(frame map[string]interface{}) {
	block.Width = frame["width"].(float64)
	block.Height = frame["height"].(float64)
	block.Left = frame["x"].(float64)
	block.Top = frame["y"].(float64)
}

func colorRGBA(collorMap map[string]interface{}) string {
	r := strconv.Itoa(int(collorMap["red"].(float64) * 255))
	g := strconv.Itoa(int(collorMap["green"].(float64) * 255))
	b := strconv.Itoa(int(collorMap["blue"].(float64) * 255))
	a := strconv.FormatFloat(collorMap["alpha"].(float64), 'f', 2, 64)
	return "rgba(" + r + ", " + g + ", " + b + ", " + a + ")"
}

func colorHex(collorMap map[string]interface{}) string {
	h := strconv.FormatInt(int64(collorMap["red"].(float64)*255), 16)
	e := strconv.FormatInt(int64(collorMap["green"].(float64)*255), 16)
	x := strconv.FormatInt(int64(collorMap["blue"].(float64)*255), 16)
	return "#" + h + e + x
}

func (s *MapShadow) boxShadow() (string, error) {
	if s.Value["isEnabled"].(bool) {
		x := strconv.Itoa(int((*s).Value["offsetX"].(float64))) + "px "
		y := strconv.Itoa(int((*s).Value["offsetY"].(float64))) + "px "
		blur := strconv.Itoa(int((*s).Value["blurRadius"].(float64))) + "px "
		color := colorRGBA((*s).Value["color"].(map[string]interface{}))
		return x + y + blur + color, nil
	}
	return "", errors.New("Disabled shadow")
}

func (block *BlockCss) fontStyle(attributedString map[string]interface{}) {
	var result Font
	attS, ok := attributedString["attributes"].([]interface{})
	if ok {
		attS, ok := attS[0].(map[string]interface{})
		if ok {
			attributes, ok := attS["attributes"].(map[string]interface{})
			if ok {
				color, ok := attributes["MSAttributedStringColorAttribute"].(map[string]interface{})
				if ok {
					result.Color = colorRGBA(color)
				}
				fontMain, ok := attributes["MSAttributedStringFontAttribute"].(map[string]interface{})
				if ok {
					font, ok := fontMain["attributes"].(map[string]interface{})
					if ok {
						result.Family = font["name"].(string)
						result.Size = font["size"].(float64)
					}
				}
				result.Spacing = attributes["kerning"].(float64)
				paragraphStyle, ok := attributes["paragraphStyle"].(map[string]interface{})
				if ok {
					lineHeight, ok := paragraphStyle["maximumLineHeight"].(float64)
					if ok {
						result.Height = lineHeight
					}
				}
			}
		}
	}
	text, ok := attributedString["string"].(string)
	if ok {
		result.Text = text
	}
	block.Font = result
}
