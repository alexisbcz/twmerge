package twmerge

import "testing"

func TestMerge(t *testing.T) {
	cases := []struct{ in, want string }{
		// From the tailwind-merge README and test suite.
		{"px-2 py-1 bg-red hover:bg-dark-red p-3 bg-[#B91C1C]", "hover:bg-dark-red p-3 bg-[#B91C1C]"},
		{"p-3 px-5", "p-3 px-5"},
		{"px-5 p-3", "p-3"},
		{"inset-x-4 right-4", "inset-x-4 right-4"},
		{"right-4 inset-x-4", "inset-x-4"},
		{"hover:bg-red-500 hover:bg-blue-500", "hover:bg-blue-500"},
		{"hover:focus:bg-red-500 focus:hover:bg-blue-500", "focus:hover:bg-blue-500"},
		{"bg-black bg-(--my-color) bg-[color:var(--x)]", "bg-[color:var(--x)]"},
		{"grid-cols-[1fr,auto] grid-cols-2", "grid-cols-2"},
		{"[mask-type:luminance] [mask-type:alpha]", "[mask-type:alpha]"},
		{"[--scroll-offset:56px] lg:[--scroll-offset:44px]", "[--scroll-offset:56px] lg:[--scroll-offset:44px]"},
		{"!p-3 !p-4 p-5", "!p-4 p-5"},
		{"p-3! p-4!", "p-4!"},
		{"-m-2 -m-5", "-m-5"},
		{"-top-12 -top-2000", "-top-2000"},
		{"text-sm text-white text-lg text-red-500", "text-lg text-red-500"},
		{"text-sm/6 text-base", "text-base"},
		{"leading-6 text-lg", "text-lg"},
		{"text-lg leading-6", "text-lg leading-6"},
		{"text-[14px] text-[#fff] text-sm", "text-[#fff] text-sm"},
		{"border-2 border-red-500 border-[3px] border-blue-500", "border-[3px] border-blue-500"},
		{"border-t-2 border-2", "border-2"},
		{"border border-input", "border border-input"},
		{"rounded-t-lg rounded-md", "rounded-md"},
		{"rounded-md rounded-t-lg", "rounded-md rounded-t-lg"},
		{"size-9 h-4 size-6", "size-6"},
		{"size-6 h-4", "size-6 h-4"},
		{"h-4 w-4 size-6", "size-6"},
		{"ring-2 ring-[3px] ring-ring/50", "ring-[3px] ring-ring/50"},
		{"outline-none outline-2 outline-dashed", "outline-2 outline-dashed"},
		{"shadow-xs shadow-lg shadow-red-500", "shadow-lg shadow-red-500"},
		{"font-medium font-bold font-mono", "font-bold font-mono"},
		{"flex block items-center flex-col flex-row flex-1", "block items-center flex-row flex-1"},
		{"justify-center justify-between justify-items-start", "justify-between justify-items-start"},
		{"overflow-x-auto overflow-hidden", "overflow-hidden"},
		{"bg-cover bg-center bg-red-500 bg-none", "bg-cover bg-center bg-red-500 bg-none"},
		{"[&_svg]:size-4 [&_svg]:size-3", "[&_svg]:size-3"},
		{"has-[>svg]:px-3 has-[>svg]:px-2.5", "has-[>svg]:px-2.5"},
		{"before:hover:p-1 hover:before:p-2", "before:hover:p-1 hover:before:p-2"},
		{"unknown-a unknown-a custom", "unknown-a unknown-a custom"},
		{"  p-1   p-2  ", "p-2"},
		{"", ""},
	}
	for _, c := range cases {
		if got := Merge(c.in); got != c.want {
			t.Errorf("Merge(%q)\n  got %q\n want %q", c.in, got, c.want)
		}
	}
}

func TestMergeMultipleLists(t *testing.T) {
	if got := Merge("h-9 px-4 rounded-md", "", "rounded-full px-2"); got != "h-9 rounded-full px-2" {
		t.Fatal(got)
	}
}

func BenchmarkMergeCached(b *testing.B) {
	in := "inline-flex items-center gap-2 rounded-md text-sm h-9 px-4 py-2 rounded-full"
	for b.Loop() {
		Merge(in)
	}
}

func BenchmarkMergeUncached(b *testing.B) {
	in := "inline-flex items-center gap-2 rounded-md text-sm font-medium whitespace-nowrap focus-visible:ring-[3px] h-9 px-4 py-2 has-[>svg]:px-3 rounded-full"
	for b.Loop() {
		merge(in)
	}
}
