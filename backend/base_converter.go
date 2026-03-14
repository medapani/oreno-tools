package backend

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type BaseConversionResult struct {
	Decimal       string `json:"decimal"`
	Hex           string `json:"hex"`
	Binary        string `json:"binary"`
	GroupedBinary string `json:"groupedBinary"`
}

type BaseCalculationResult struct {
	Decimal       string `json:"decimal"`
	Hex           string `json:"hex"`
	Binary        string `json:"binary"`
	GroupedBinary string `json:"groupedBinary"`
}

func ConvertBaseValue(input string, base int, bitWidth int, signed bool) (BaseConversionResult, error) {
	bits, err := bitWidthToUint(bitWidth)
	if err != nil {
		return BaseConversionResult{}, err
	}
	parsed, err := parseBigIntByBase(input, base)
	if err != nil {
		return BaseConversionResult{}, err
	}

	clamped := clampBigInt(parsed, bits, signed)
	binary := formatBinaryForBitWidth(clamped, bitWidth, bits)

	return BaseConversionResult{
		Decimal:       clamped.Text(10),
		Hex:           formatHexForBitWidth(clamped, bitWidth, bits),
		Binary:        binary,
		GroupedBinary: groupBinary(binary),
	}, nil
}

func CalculateBaseExpression(aInput string, aBase int, operator string, bInput string, bBase int) (BaseCalculationResult, error) {
	a, err := parseBigIntByBase(aInput, aBase)
	if err != nil {
		return BaseCalculationResult{}, fmt.Errorf("invalid A: %w", err)
	}
	b, err := parseBigIntByBase(bInput, bBase)
	if err != nil {
		return BaseCalculationResult{}, fmt.Errorf("invalid B: %w", err)
	}

	result, err := executeBigIntOperation(a, b, operator)
	if err != nil {
		return BaseCalculationResult{}, err
	}

	binary := formatRawBinary(result)
	return BaseCalculationResult{
		Decimal:       result.Text(10),
		Hex:           formatRawHex(result),
		Binary:        binary,
		GroupedBinary: groupBinaryWithSign(binary),
	}, nil
}

func parseBigIntByBase(input string, base int) (*big.Int, error) {
	if base != 2 && base != 10 && base != 16 {
		return nil, fmt.Errorf("unsupported base: %d", base)
	}
	trimmed := strings.TrimSpace(input)
	if trimmed == "" || trimmed == "-" {
		return nil, fmt.Errorf("empty value")
	}
	value, ok := new(big.Int).SetString(trimmed, base)
	if !ok {
		return nil, fmt.Errorf("cannot parse value")
	}
	return value, nil
}

func clampBigInt(n *big.Int, bits uint, signed bool) *big.Int {
	v := new(big.Int).Set(n)

	if signed {
		limit := new(big.Int).Lsh(big.NewInt(1), bits-1)
		upperBound := new(big.Int).Sub(limit, big.NewInt(1))
		lowerBound := new(big.Int).Neg(limit)
		if v.Cmp(lowerBound) < 0 {
			return lowerBound
		}
		if v.Cmp(upperBound) > 0 {
			return upperBound
		}
		return v
	}

	maxValue := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits), big.NewInt(1))
	if v.Sign() < 0 {
		return big.NewInt(0)
	}
	if v.Cmp(maxValue) > 0 {
		return maxValue
	}
	return v
}

func formatBinaryForBitWidth(n *big.Int, bitWidth int, bits uint) string {
	if n.Sign() < 0 {
		mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits), big.NewInt(1))
		twos := new(big.Int).And(n, mask).Text(2)
		if len(twos) < bitWidth {
			return strings.Repeat("0", bitWidth-len(twos)) + twos
		}
		return twos
	}
	bin := n.Text(2)
	if len(bin) < bitWidth {
		return strings.Repeat("0", bitWidth-len(bin)) + bin
	}
	return bin
}

func formatHexForBitWidth(n *big.Int, bitWidth int, bits uint) string {
	if n.Sign() < 0 {
		mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), bits), big.NewInt(1))
		twos := strings.ToUpper(new(big.Int).And(n, mask).Text(16))
		hexLen := bitWidth / 4
		if len(twos) < hexLen {
			return strings.Repeat("0", hexLen-len(twos)) + twos
		}
		return twos
	}
	return strings.ToUpper(n.Text(16))
}

func formatRawHex(n *big.Int) string {
	if n.Sign() < 0 {
		abs := new(big.Int).Abs(new(big.Int).Set(n))
		return "-" + strings.ToUpper(abs.Text(16))
	}
	return strings.ToUpper(n.Text(16))
}

func formatRawBinary(n *big.Int) string {
	if n.Sign() < 0 {
		abs := new(big.Int).Abs(new(big.Int).Set(n))
		return "-" + abs.Text(2)
	}
	return n.Text(2)
}

func groupBinary(bin string) string {
	if len(bin) == 0 {
		return ""
	}
	parts := make([]string, 0, (len(bin)+3)/4)
	for i := len(bin); i > 0; i -= 4 {
		start := i - 4
		if start < 0 {
			start = 0
		}
		parts = append(parts, bin[start:i])
	}
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, " ")
}

func groupBinaryWithSign(bin string) string {
	if strings.HasPrefix(bin, "-") {
		return "-" + groupBinary(bin[1:])
	}
	return groupBinary(bin)
}

func executeBigIntOperation(a *big.Int, b *big.Int, operator string) (*big.Int, error) {
	switch operator {
	case "+":
		return new(big.Int).Add(a, b), nil
	case "-":
		return new(big.Int).Sub(a, b), nil
	case "*":
		return new(big.Int).Mul(a, b), nil
	case "/":
		if b.Sign() == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return new(big.Int).Quo(a, b), nil
	case "%":
		if b.Sign() == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		return new(big.Int).Rem(a, b), nil
	case "&":
		return new(big.Int).And(a, b), nil
	case "|":
		return new(big.Int).Or(a, b), nil
	case "^":
		return new(big.Int).Xor(a, b), nil
	case "<<":
		if b.Sign() < 0 {
			return nil, fmt.Errorf("shift count must be non-negative")
		}
		if b.BitLen() > 64 {
			return nil, fmt.Errorf("shift count is too large")
		}
		shift := b.Uint64()
		if shift > 1<<30 {
			return nil, fmt.Errorf("shift count is too large")
		}
		maxUint := ^uint(0)
		if strconv.IntSize == 32 && shift > uint64(maxUint) {
			return nil, fmt.Errorf("shift count is too large")
		}
		return new(big.Int).Lsh(a, uint(shift)), nil
	case ">>":
		if b.Sign() < 0 {
			return nil, fmt.Errorf("shift count must be non-negative")
		}
		if b.BitLen() > 64 {
			return nil, fmt.Errorf("shift count is too large")
		}
		shift := b.Uint64()
		if shift > 1<<30 {
			return nil, fmt.Errorf("shift count is too large")
		}
		maxUint := ^uint(0)
		if strconv.IntSize == 32 && shift > uint64(maxUint) {
			return nil, fmt.Errorf("shift count is too large")
		}
		return new(big.Int).Rsh(a, uint(shift)), nil
	default:
		return nil, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func bitWidthToUint(bitWidth int) (uint, error) {
	switch bitWidth {
	case 8:
		return 8, nil
	case 16:
		return 16, nil
	case 32:
		return 32, nil
	case 64:
		return 64, nil
	default:
		return 0, fmt.Errorf("unsupported bit width: %d", bitWidth)
	}
}
