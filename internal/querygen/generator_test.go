package querygen

import "testing"

func TestClassifySourceType(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		keywords string
		want     string
	}{
		{name: "mall", source: "静安大悦城", keywords: "商场,快闪,打卡有礼", want: "mall"},
		{name: "brand", source: "兰蔻LANCOME", keywords: "品牌,新品,赠礼", want: "brand"},
		{name: "event", source: "史努比75周年", keywords: "官方活动,巡展,联名", want: "official_event"},
		{name: "info", source: "上海情报站", keywords: "探店,情报,汇总", want: "info_account"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ClassifySourceType(tc.source, tc.keywords)
			if got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
}

func TestGenerateQueries(t *testing.T) {
	source := Source{ID: 1, Name: "环贸iapm商场", Keywords: "iapm,活动日历,快闪,打卡有礼,会员礼", SourceType: "mall"}
	queries := GenerateQueries(source, 3)
	if len(queries) != 3 {
		t.Fatalf("got %d queries", len(queries))
	}
	want := []string{"环贸iapm商场 活动日历", "环贸iapm商场 快闪", "环贸iapm商场 打卡有礼"}
	for i, q := range want {
		if queries[i] != q {
			t.Fatalf("query[%d] got %s want %s", i, queries[i], q)
		}
	}
}
