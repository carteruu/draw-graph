package draw

import (
	"drawgraph"
	_ "embed"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/shomali11/gridder"
	"image/color"
	"log"
	"os"
	"unicode/utf8"
)

type point struct {
	x, y int
}

//go:embed Alimama_ShuHeiTi_Bold.ttf
var ttcBytes []byte

const (
	nodeHeight     = 6  //节点高度-格子数
	plaidPixel     = 10 //格子边长-像素
	fontSize       = 12 //字体大小
	textWidthPixel = 15 //文字的宽度-像素
)

var (
	logger   = log.Default()
	ff, _    = truetype.Parse(ttcBytes)
	fontFace = truetype.NewFace(ff, &truetype.Options{Size: float64(fontSize)})

	stringConfigBlack = gridder.StringConfig{Color: color.Black}
	stringConfigGray  = gridder.StringConfig{Color: color.Gray{Y: 0x000080}}

	pathConfig = gridder.PathConfig{
		Color: color.Gray{Y: 220},
	}
)

// Draw 画图
func Draw(name, outDir string, nodes []draw_graph.Node) error {
	nodeMap := make(map[int]draw_graph.Node, len(nodes))
	nodeChildren := make(map[int][]int, len(nodes))
	for _, n := range nodes {
		if n.ID == 0 {
			return fmt.Errorf("id can not 0")
		}
		if _, ok := nodeMap[n.ID]; ok {
			return fmt.Errorf("duplicate node: %v", n.ID)
		}
		nodeMap[n.ID] = n
		nodeChildren[n.Parent] = append(nodeChildren[n.Parent], n.ID)
	}
	return draw(name, outDir, nodeMap, nodeChildren)
}

// 需要画的节点
func needDrawNode(n int) bool {
	return true
}

// 找到需要画的子节点（跳过不需要画的节点）
func findNeedDrawChildren(childrenMap map[int][]int, n int) (children []int) {
	if !needDrawNode(n) {
		return nil
	}
	for _, child := range childrenMap[n] {
		if needDrawNode(child) {
			children = append(children, child)
			continue
		}
		children = append(children, findNeedDrawChildren(childrenMap, child)...)
	}
	return children
}

// 画图
func draw(name, outDir string, nodeMap map[int]draw_graph.Node, nodeChildren map[int][]int) error {
	//跳过不需要画的节点，重新生成边
	childrenSet := make(map[int]map[int]struct{}, len(nodeChildren))
	for id := range nodeChildren {
		if !needDrawNode(id) {
			continue
		}
		childrenSet[id] = make(map[int]struct{})
		for _, child := range findNeedDrawChildren(nodeChildren, id) {
			childrenSet[id][child] = struct{}{}
		}
	}
	//计算图纸所需的长宽
	rowNum, colNodeNumMax, textLenMax, err := drawInfo(childrenSet, nodeMap)
	if err != nil {
		return err
	}
	//节点宽度-格子数
	nodeWidth := textLenMax * textWidthPixel / plaidPixel
	gridConfig := gridder.GridConfig{
		Columns: nodeWidth * (colNodeNumMax + 1),
		Rows:    nodeHeight * rowNum,
	}
	imageConfig := gridder.ImageConfig{
		Width:  gridConfig.Columns * plaidPixel,
		Height: gridConfig.Rows * plaidPixel,
		Name:   fmt.Sprintf("%s.png", name),
	}
	grid, err := gridder.New(imageConfig, gridConfig)
	if err != nil {
		return err
	}

	parentPointMap := make(map[int][]point)
	y := nodeHeight / 2
	inMap := inCnt(childrenSet)
	q := make([]int, 0, len(childrenSet))
	for id, in := range inMap {
		if in == 0 {
			q = append(q, id)
		}
	}
	for len(q) > 0 {
		//节点对称排列：一层奇数个节点时，中间的节点居中；偶数个节点时，分布在左右两边
		size := len(q)
		offset := 0
		if size&1 == 0 {
			offset = nodeWidth / 2
		}
		for i := 0; i < size; i++ {
			id := q[0]
			q = q[1:]
			stringConfig := stringConfigBlack
			//当前节点的横坐标
			x := gridConfig.Columns/2 + (i-size/2)*nodeWidth + offset
			if err := grid.DrawString(y, x, fmt.Sprintf("%v - %v", id, nodeMap[id].Text), fontFace, stringConfig); err != nil {
				return err
			}
			//当前节点到父节点的边
			for _, parentPoint := range parentPointMap[id] {
				currentX := x
				if x < parentPoint.x {
					parentPoint.x--
					currentX++
				} else if x > parentPoint.x {
					parentPoint.x++
					currentX--
				}
				parentY := parentPoint.y
				//当前节点row-1是为了让边上移一点，父节点row-1是为了让边下移一点，避免边遮住文字
				if err := grid.DrawPath(y-1, currentX, parentY+1, parentPoint.x, pathConfig); err != nil {
					return err
				}
			}
			//子节点
			for child := range childrenSet[id] {
				inMap[child]--
				parentPointMap[child] = append(parentPointMap[child], point{x: x, y: y})
				if inMap[child] == 0 {
					q = append(q, child)
				}
			}
		}
		y += nodeHeight
	}
	filePath := outDir + "/" + imageConfig.Name
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if err := grid.EncodePNG(file); err != nil {
		return err
	}
	logger.Printf("生成图片: %s\n", filePath)
	return nil
}

func inCnt(childrenSet map[int]map[int]struct{}) map[int]int {
	inMap := make(map[int]int, len(childrenSet))
	for parent, cs := range childrenSet {
		if parent == 0 {
			continue
		}
		if _, ok := inMap[parent]; !ok {
			inMap[parent] = 0
		}
		for id := range cs {
			inMap[id]++
		}
	}
	return inMap
}

// 计算图纸的行数、每行最大节点数，文本最大字符数
func drawInfo(childrenSet map[int]map[int]struct{}, nodeMap map[int]draw_graph.Node) (int, int, int, error) {
	inMap := inCnt(childrenSet)
	//对话最大字符数
	textLenMax := 10
	//每行最大节点数
	colNodeNumMax := 0
	//行数
	rowNum := 1
	q := make([]int, 0, len(childrenSet))
	for id, in := range inMap {
		if in == 0 {
			q = append(q, id)
		}
	}
	visitedIDs := make(map[int]struct{})
	for len(q) > 0 {
		rowNum++
		size := len(q)
		colNodeNumMax = draw_graph.MaxInt(colNodeNumMax, size)
		for i := 0; i < size; i++ {
			n := q[0]
			visitedIDs[n] = struct{}{}
			textLenMax = draw_graph.MaxInt(textLenMax, utf8.RuneCountInString(nodeMap[n].Text))
			q = q[1:]
			for child := range childrenSet[n] {
				inMap[child]--
				if inMap[child] == 0 {
					q = append(q, child)
				}
			}
		}
	}
	if len(visitedIDs) != len(nodeMap) {
		ns := make([]int, 0, len(nodeMap)-len(visitedIDs))
		for id := range nodeMap {
			if _, ok := visitedIDs[id]; !ok {
				ns = append(ns, id)
			}
		}
		return 0, 0, 0, fmt.Errorf("存在环:%v", ns)
	}
	return rowNum, colNodeNumMax, textLenMax, nil
}
