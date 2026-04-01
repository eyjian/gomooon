// Package mooonstr
// Wrote by yijian on 2025/06/20
package mooonstr

import "testing"

// go test -v -run="TestCompressChinese"
func TestCompressChinese(t *testing.T) {
	s := `我是中国人，我人爱中国。中国具有5000年文明史，是一个伟大的国家。
三皇五帝始，尧舜禹相传；
夏商与西周，东周分两段；
春秋和战国，一统秦两汉；
三分魏蜀吴，二晋前后延；
南北朝并立，隋唐五代传；
宋元明清后，皇朝至此完。
五代十国：
五代：梁唐晋汉周，前边都有后；
十国：前后蜀，南北汉，南唐南平曾为伴，
吴越吴闽楚十国，割据混战中原乱。
南北朝：
南朝：宋齐梁陈相交替；
北朝：北魏分东西（东魏、西魏），北周灭北齐。
五胡十六国：
前后南三燕，西秦南凉鲜卑建；
前西二凉和北燕，汉族政权仍延续；
前赵北凉夏匈奴，前秦后凉汉（成汉）氐建；
羯后赵，羌后秦，十六小国长混战。
`
	compressed, err := CompressChinese(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d => %d\n", len(s), len(compressed))

	s, err = DecompressChinese(compressed)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", s)
}
