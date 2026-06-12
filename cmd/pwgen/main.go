package main

import (
	"flag"
	"fmt"
	"os"

	"password_gen/internal/generator"
)

func main() {
	var length int
	flag.IntVar(&length, "l", 16, "password length (6-128)")
	flag.IntVar(&length, "length", 16, "password length (6-128)")

	var count int
	flag.IntVar(&count, "c", 1, "number of passwords to generate")
	flag.IntVar(&count, "count", 1, "number of passwords to generate")

	var output string
	flag.StringVar(&output, "o", "", "output file (default stdout)")
	flag.StringVar(&output, "output", "", "output file (default stdout)")

	var noUpper, noLower, noDigit, noSymbol bool
	flag.BoolVar(&noUpper, "no-upper", false, "exclude uppercase letters")
	flag.BoolVar(&noLower, "no-lower", false, "exclude lowercase letters")
	flag.BoolVar(&noDigit, "no-digit", false, "exclude digits")
	flag.BoolVar(&noSymbol, "no-symbol", false, "exclude symbols")

	var showStrength bool
	flag.BoolVar(&showStrength, "s", false, "show password strength")
	flag.BoolVar(&showStrength, "strength", false, "show password strength")

	flag.Parse()

	if length < 6 || length > 128 {
		fmt.Printf("length must be between 6 and 128\n")
		os.Exit(1)
	}
	if count < 1 || count > 1000 {
		fmt.Printf("count must be between 1 and 1000\n")
		os.Exit(1)
	}

	// 默认启用所有字符集，被 no-xxx 关闭则停用
	opts := generator.Options{
		Length:    length,
		UseUpper:  !noUpper,
		UseLower:  !noLower,
		UseDigit:  !noDigit,
		UseSymbol: !noSymbol,
	}

	// 至少要启用一类
	if !opts.UseUpper && !opts.UseLower && !opts.UseDigit && !opts.UseSymbol {
		fmt.Printf("at least one character class must be enabled\n")
		os.Exit(1)
	}

	g := generator.New()

	// 计算字符集大小用于熵计算
	charsetSize := 0
	if opts.UseLower {
		charsetSize += len(generator.CharsetLower)
	}
	if opts.UseUpper {
		charsetSize += len(generator.CharsetUpper)
	}
	if opts.UseDigit {
		charsetSize += len(generator.CharsetDigit)
	}
	if opts.UseSymbol {
		charsetSize += len(generator.CharsetSymbol)
	}

	// 准备输出
	var writer *os.File
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			fmt.Printf("cannot create output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		writer = f
	} else {
		writer = os.Stdout
	}

	// 批量生成
	for i := 0; i < count; i++ {
		password, err := g.Generate(opts)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if showStrength {
			entropy := generator.Entropy(length, charsetSize)
			fmt.Fprintf(writer, "%s\t%.1f bits\n", password, entropy)
		} else {
			fmt.Fprintln(writer, password)
		}
	}
}