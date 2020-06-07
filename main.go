package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var big1, bigA, bigN = big.NewInt(1), big.NewInt(0), big.NewInt(0)

// factorial tail recursive factorial function.
func factorial(x uint64) *big.Int {
	if x <= 1 {
		return big1
	}
	var F func(n uint64, a *big.Int) *big.Int
	F = func(n uint64, a *big.Int) *big.Int {
		if n <= 1 {
			return a
		}
		bigN.SetInt64(int64(n))
		bigA.Mul(a, bigN)
		return F(n-1, bigA)
	}
	return F(x, big1)
}

var nbig, kbig, nmkbig, resbig = big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)

// nchoosek big.Int nchoosek.
func nchoosek(n, k uint64) *big.Int {
	if k >= n {
		return big1
	}
	kbig.Set(factorial(k))
	nbig.Set(factorial(n))
	nmkbig.Set(factorial(n - k))

	resbig.Mul(kbig, nmkbig)
	resbig.Div(nbig, resbig)

	return resbig
}

// Pascal main struct.
type Pascal struct {
	Depth uint64
	T     [][]*big.Int
}

// NewPascal returns a new Pascal triangle of depth n.
func NewPascal(n uint64, only bool) *Pascal {
	var p Pascal
	if only {
		p = Pascal{Depth: 1, T: make([][]*big.Int, 1, 1)}
	} else {
		p = Pascal{Depth: n, T: make([][]*big.Int, n+1, n+1)}
	}
	var i, j uint64
	if only {
		p.T[0] = make([]*big.Int, n+1, n+1)
		for j = 0; j <= n; j++ {
			var bigi = big.NewInt(0)
			bigi.Set(nchoosek(n, j))
			p.T[0][j] = bigi
		}
		return &p
	}
	for i = 0; i <= n; i++ {
		p.T[i] = make([]*big.Int, i+1, i+1)
		for j = 0; j <= i; j++ {
			var bigi = big.NewInt(0)
			bigi.Set(nchoosek(i, j))
			p.T[i][j] = bigi
		}
	}

	return &p
}

// Max returns the maximum value in the triangle: nchoosek(depth, depth/2)
func (p *Pascal) Max() *big.Int {
	return nchoosek(p.Depth, p.Depth/2)
}

// ToText uses tabwriter.Writer to pretty print to text.
func (p *Pascal) ToText(rowHeaders bool) string {
	var b strings.Builder
	if width == 0 {
		width = len(p.Max().String()) + 1
	}
	w := tabwriter.NewWriter(&b, width, 0, 0, ' ', tabwriter.AlignRight)

	for i, row := range p.T {
		if rowHeaders {
			fmt.Fprintf(w, "%d:\t", i)
		}
		var jj = int((2 * p.Depth) + 1)
		fmt.Fprint(w, strings.Repeat("\t", int(p.Depth)-i))
		for _, j := range row {
			fmt.Fprintf(w, "%s\t\t", j.String())
		}
		fmt.Fprintf(w, "%s\n", strings.Repeat("\t", jj-(int(p.Depth)-i)))
	}
	w.Flush()

	return b.String()
}

// ToHTML generates an HTML table with the values.
func (p *Pascal) ToHTML(rowHeaders bool) string {
	var b strings.Builder
	max := p.Max()
	fmt.Fprintf(&b, `
<!DOCTYPE html>
<html>
<body>
<h1>Pascal triangle of depth %d</h1>
<nav>
	<ul>
		<li><a href="#central-column">Central column</a></li>
		<li><a href="#biggest-number">Biggest number</a></li>
	</ul>
</nav>
<p>Biggest number is %s.</p>
<p>It has %d digits.</p>
`, p.Depth, max.String(), len(max.String()))
	b.WriteString(`<style>
	table {
		empty-cells: show;
	}
	table, td, th {
		text-align: right;
		font-size: 9px;
	}
</style>
`)
	b.WriteString("<table><tbody>\n")

	for i, row := range p.T {
		b.WriteString("<tr>")
		if rowHeaders {
			fmt.Fprintf(&b, `<th scope="row">%d</th>`, i)
		}
		var jj = int((2 * p.Depth) + 1)
		fmt.Fprint(&b, strings.Repeat(`<td></td>`, int(p.Depth)-i))
		for _, j := range row {
			if i == 0 {
				fmt.Fprintf(&b, `<td id="central-column">%s</td><td></td>`, j.String())
			} else if j.Cmp(max) == 0 {
				fmt.Fprintf(&b, `<td id="biggest-number" style="color: red;">%s</td><td></td>`, j.String())
			} else {
				fmt.Fprintf(&b, "<td>%s</td><td></td>", j.String())
			}
		}
		fmt.Fprintf(&b, "%s</tr>\n", strings.Repeat("<td></td>", jj-(int(p.Depth)-i)))
	}
	b.WriteString(`</tbody>
</table>
</body>
</html>`)

	return b.String()
}

var (
	headers bool
	width   int
	format  string
	only    bool
	fact    bool
	biggest bool
	choose  bool
)

func init() {
	const (
		defaultHeaders = false
		usageHeaders   = "print row headers in triangle"
		defaultWidth   = 0
		usageWidth     = "min width for each \"cell\" in text format, (default len(nchoosek(depth, depth/2))+1)"
		defaultFormat  = "text"
		usageFormat    = "Format to output. Options are text, html, raw"
		defaultOnly    = false
		usageOnly      = "only output the row at 'depth'"
		defaultFact    = false
		usageFact      = "calculate the factorial of the argument"
		defaultBiggest = false
		usageBiggest   = "return the maximum value in the triangle: nchoosek(depth, depth/2)"
		defaultChoose  = false
		usageChoose    = "accept a second argument k and calculate nchoosek(depth, k)"
	)
	flag.BoolVar(&headers, "headers", defaultHeaders, usageHeaders)
	flag.BoolVar(&headers, "n", defaultHeaders, usageHeaders+"(shorthand)")
	flag.IntVar(&width, "width", defaultWidth, usageWidth)
	flag.IntVar(&width, "w", defaultWidth, usageWidth+"(shorthand)")
	flag.StringVar(&format, "format", defaultFormat, usageFormat)
	flag.StringVar(&format, "f", defaultFormat, usageFormat+"(shorthand)")
	flag.BoolVar(&only, "only", defaultOnly, usageOnly)
	flag.BoolVar(&only, "o", defaultOnly, usageOnly+"(shorthand)")
	flag.BoolVar(&fact, "factorial", defaultFact, usageFact)
	flag.BoolVar(&fact, "y", defaultFact, usageFact)
	flag.BoolVar(&biggest, "biggest", defaultBiggest, usageBiggest)
	flag.BoolVar(&biggest, "b", defaultBiggest, usageBiggest+"(shorthand)")
	flag.BoolVar(&choose, "choosek", defaultChoose, usageChoose)
	flag.BoolVar(&choose, "c", defaultChoose, usageChoose+"(shorthand)")
	flag.Usage = func() {
		usageStr := fmt.Sprintf(`
Usage: %s [OPTIONS] depth [k]

-f, -format string	%s (default: "%s")
-n, -headers		%s
-o, -only		%s
-y, -factorial		%s
-b, -biggest		%s
-w, -width		%s
-c, -choosek	%s
`, os.Args[0], usageFormat, defaultFormat, usageHeaders, usageOnly, usageFact, usageBiggest, usageWidth, usageChoose)
		fmt.Fprintf(flag.CommandLine.Output(), "%s", usageStr)
	}
}

func main() {
	flag.Parse()
	v, err := strconv.ParseUint(flag.Arg(0), 10, 64)
	if err != nil {
		panic(err)
	}
	if fact {
		fmt.Println(factorial(v))
		os.Exit(0)
	}
	if biggest {
		fmt.Println(nchoosek(v, v/2))
		os.Exit(0)
	}
	if choose {
		k, err := strconv.ParseUint(flag.Arg(1), 10, 64)
		if err != nil {
			panic(err)
		}
		fmt.Println(nchoosek(v, k))
		os.Exit(0)
	}
	p := NewPascal(v, only)
	switch format {
	case "html":
		fmt.Fprintf(os.Stdout, "%s", p.ToHTML(headers))
	case "raw":
		for i, row := range p.T {
			if headers {
				fmt.Fprintf(os.Stdout, "%d:\t", i)
			}
			fmt.Fprintf(os.Stdout, "%v\n", row)
		}
	default:
		fmt.Fprintf(os.Stdout, "%s", p.ToText(headers))
	}
}
