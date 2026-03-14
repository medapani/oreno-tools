import React, { useState, useCallback, useEffect } from 'react';
import { CalculateBaseExpression, ConvertBaseValue } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

type Operator = '+' | '-' | '*' | '/' | '%' | '&' | '|' | '^' | '<<' | '>>';

type BitWidth = 8 | 16 | 32 | 64;

const BIT_WIDTHS: BitWidth[] = [8, 16, 32, 64];

function groupBinary(bin: string): string {
  // 4ビットごとにスペース区切り (右から)
  const reversed = bin.split('').reverse();
  const groups: string[] = [];
  for (let i = 0; i < reversed.length; i += 4) {
    groups.push(reversed.slice(i, i + 4).reverse().join(''));
  }
  return groups.reverse().join(' ');
}

const OPERATORS: { label: string; value: Operator; title: string }[] = [
  { label: '+', value: '+', title: '加算' },
  { label: '-', value: '-', title: '減算' },
  { label: '×', value: '*', title: '乗算' },
  { label: '÷', value: '/', title: '除算' },
  { label: '%', value: '%', title: '剰余' },
  { label: '&', value: '&', title: 'AND' },
  { label: '|', value: '|', title: 'OR' },
  { label: '^', value: '^', title: 'XOR' },
  { label: '<<', value: '<<', title: '左シフト' },
  { label: '>>', value: '>>', title: '右シフト' },
];

type CalcBase = 2 | 10 | 16;

function calcInputPlaceholder(base: CalcBase): string {
  if (base === 2) return '例: 1010';
  if (base === 16) return '例: FF';
  return '例: 255';
}

function calcInputPattern(base: CalcBase): RegExp {
  if (base === 2) return /^-?[01]*$/;
  if (base === 16) return /^-?[0-9a-fA-F]*$/;
  return /^-?\d*$/;
}

export default function BaseConverter() {
  const [binInput, setBinInput] = useSessionState('baseconv.bin', '');
  const [decInput, setDecInput] = useSessionState('baseconv.dec', '');
  const [hexInput, setHexInput] = useSessionState('baseconv.hex', '');
  const [bitWidth, setBitWidth] = useSessionState<BitWidth>('baseconv.bitWidth', 32);
  const [signed, setSigned] = useSessionState('baseconv.signed', true);
  const [error, setError] = useState('');
  const [binaryGrouped, setBinaryGrouped] = useState('');

  // 四則演算
  const [calcA, setCalcA] = useSessionState('baseconv.calcA', '');
  const [calcB, setCalcB] = useSessionState('baseconv.calcB', '');
  const [calcBaseA, setCalcBaseA] = useSessionState<CalcBase>('baseconv.calcBaseA', 10);
  const [calcBaseB, setCalcBaseB] = useSessionState<CalcBase>('baseconv.calcBaseB', 10);
  const [calcOp, setCalcOp] = useSessionState<Operator>('baseconv.calcOp', '+');
  const [calcResult, setCalcResult] = useState<{ dec: string; hex: string; bin: string; grouped: string } | null>(null);
  const [calcError, setCalcError] = useState('');

  const handleCalcAChange = (v: string) => {
    const upper = calcBaseA === 16 ? v.toUpperCase() : v;
    setCalcA(upper);
  };

  const handleCalcBChange = (v: string) => {
    const upper = calcBaseB === 16 ? v.toUpperCase() : v;
    setCalcB(upper);
  };

  const compute = useCallback(async () => {
    setCalcError('');
    setCalcResult(null);
    if (calcA.trim() === '' || calcA.trim() === '-' || calcB.trim() === '' || calcB.trim() === '-') {
      return;
    }
    if (!calcInputPattern(calcBaseA).test(calcA.trim())) {
      setCalcError('A の値が不正です');
      return;
    }
    if (!calcInputPattern(calcBaseB).test(calcB.trim())) {
      setCalcError('B の値が不正です');
      return;
    }
    try {
      const result = await CalculateBaseExpression(calcA, calcBaseA, calcOp, calcB, calcBaseB);
      setCalcResult({
        dec: result.decimal,
        hex: result.hex,
        bin: result.binary,
        grouped: result.groupedBinary,
      });
    } catch (e) {
      const message = String(e);
      if (message.includes('division by zero')) {
        setCalcError('0 での除算はできません');
        return;
      }
      if (message.includes('modulo by zero')) {
        setCalcError('0 での剰余はできません');
        return;
      }
      if (message.includes('shift count must be non-negative')) {
        setCalcError('シフト量は 0 以上の値を指定してください');
        return;
      }
      setCalcError(`計算に失敗しました: ${message}`);
    }
  }, [calcA, calcB, calcBaseA, calcBaseB, calcOp]);

  useEffect(() => {
    void compute();
  }, [compute]);

  const convertFrom = useCallback(
    async (value: string, base: CalcBase, source: 'bin' | 'dec' | 'hex') => {
      setError('');
      if (value.trim() === '' || value.trim() === '-') {
        if (source !== 'bin') setBinInput('');
        if (source !== 'dec') setDecInput('');
        if (source !== 'hex') setHexInput('');
        setBinaryGrouped('');
        return;
      }

      if (base === 2 && !/^-?[01]*$/.test(value.trim())) {
        setError('2進数は 0 と 1 のみ入力できます');
        return;
      }
      if (base === 10 && !/^-?\d*$/.test(value.trim())) {
        setError('10進数は数字のみ入力できます');
        return;
      }
      if (base === 16 && !/^-?[0-9a-fA-F]*$/.test(value.trim())) {
        setError('16進数は 0-9, A-F のみ入力できます');
        return;
      }

      try {
        const result = await ConvertBaseValue(value, base, bitWidth, signed);
        if (source !== 'bin') setBinInput(result.binary);
        if (source !== 'dec') setDecInput(result.decimal);
        if (source !== 'hex') setHexInput(result.hex);
        setBinaryGrouped(result.groupedBinary);
      } catch (e) {
        setError(`変換に失敗しました: ${String(e)}`);
      }
    },
    [bitWidth, signed, setBinInput, setDecInput, setHexInput],
  );

  const handleBinChange = (value: string) => {
    setBinInput(value);
    void convertFrom(value, 2, 'bin');
  };

  const handleDecChange = (value: string) => {
    setDecInput(value);
    void convertFrom(value, 10, 'dec');
  };

  const handleHexChange = (value: string) => {
    const upper = value.toUpperCase();
    setHexInput(upper);
    void convertFrom(upper, 16, 'hex');
  };

  const handleClear = () => {
    setBinInput('');
    setDecInput('');
    setHexInput('');
    setBinaryGrouped('');
    setError('');
  };

  // ビット幅や符号変更時は入力欄を上書きしない
  const handleBitWidthChange = (newWidth: BitWidth) => {
    setBitWidth(newWidth);
  };

  const handleSignedChange = (isSigned: boolean) => {
    setSigned(isSigned);
  };

  // 入力値は保持したまま、表示用2進数のみビット幅/符号に応じて更新する
  useEffect(() => {
    const refreshPreview = async () => {
      if (decInput.trim() === '' || decInput.trim() === '-' || !/^-?\d+$/.test(decInput.trim())) {
        setBinaryGrouped('');
        return;
      }
      try {
        const result = await ConvertBaseValue(decInput, 10, bitWidth, signed);
        setBinaryGrouped(result.groupedBinary);
      } catch {
        setBinaryGrouped('');
      }
    };

    void refreshPreview();
  }, [decInput, bitWidth, signed]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-2xl">
        <h1 className="text-3xl font-bold mb-2 text-center">基数変換/計算</h1>
        <p className="text-gray-400 text-sm text-center mb-6">2進数 / 10進数 / 16進数 の相互変換</p>

        {/* ビット幅 & 符号 */}
        <div className="flex flex-wrap gap-4 mb-6 items-center">
          <div>
            <span className="text-sm font-semibold text-gray-300 mr-2">ビット幅:</span>
            {BIT_WIDTHS.map((w) => (
              <button
                key={w}
                type="button"
                onClick={() => handleBitWidthChange(w)}
                className={`mr-2 px-3 py-1 rounded text-sm font-semibold transition ${bitWidth === w ? 'bg-blue-600 text-white' : 'bg-gray-700 hover:bg-gray-600 text-gray-200'
                  }`}
              >
                {w}bit
              </button>
            ))}
          </div>
          <div className="flex items-center gap-2">
            <span className="text-sm font-semibold text-gray-300">符号:</span>
            <button
              type="button"
              onClick={() => handleSignedChange(true)}
              className={`px-3 py-1 rounded text-sm font-semibold transition ${signed ? 'bg-blue-600 text-white' : 'bg-gray-700 hover:bg-gray-600 text-gray-200'
                }`}
            >
              符号あり
            </button>
            <button
              type="button"
              onClick={() => handleSignedChange(false)}
              className={`px-3 py-1 rounded text-sm font-semibold transition ${!signed ? 'bg-blue-600 text-white' : 'bg-gray-700 hover:bg-gray-600 text-gray-200'
                }`}
            >
              符号なし
            </button>
          </div>
        </div>

        {/* 入力フィールド */}
        <div className="space-y-4 mb-6">
          {/* 2進数 */}
          <div>
            <label className="block text-sm font-semibold mb-1 text-gray-300">
              2進数 (Binary)
            </label>
            <input
              type="text"
              value={binInput}
              onChange={(e) => handleBinChange(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono"
              placeholder="例: 1010"
              spellCheck={false}
            />
          </div>

          {/* 10進数 */}
          <div>
            <label className="block text-sm font-semibold mb-1 text-gray-300">
              10進数 (Decimal)
            </label>
            <input
              type="text"
              value={decInput}
              onChange={(e) => handleDecChange(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono"
              placeholder="例: 255"
              spellCheck={false}
            />
          </div>

          {/* 16進数 */}
          <div>
            <label className="block text-sm font-semibold mb-1 text-gray-300">
              16進数 (Hexadecimal)
            </label>
            <div className="relative">
              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 font-mono select-none">0x</span>
              <input
                type="text"
                value={hexInput}
                onChange={(e) => handleHexChange(e.target.value)}
                className="w-full p-3 pl-9 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono uppercase"
                placeholder="例: FF"
                spellCheck={false}
              />
            </div>
          </div>
        </div>

        {/* エラー */}
        {error && (
          <div className="mb-4 p-3 rounded bg-red-900 border border-red-700 text-red-300 text-sm">
            {error}
          </div>
        )}

        {/* クリアボタン */}
        <div className="mb-6">
          <button
            type="button"
            onClick={handleClear}
            className="px-4 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
          >
            クリア
          </button>
        </div>

        {/* ビット表示 */}
        {binaryGrouped && (
          <div className="bg-gray-700 rounded-lg p-4">
            <div className="text-xs font-semibold text-gray-400 mb-2">
              {bitWidth}bit 2進数表現 (4bit区切り)
            </div>
            <div className="font-mono text-green-400 text-lg break-all tracking-wider">
              {binaryGrouped}
            </div>
          </div>
        )}

        {/* 範囲情報 */}
        <div className="mt-4 mb-8 text-xs text-gray-500 space-y-1">
          <div>
            表現範囲:{' '}
            {signed
              ? `${(-(1n << (BigInt(bitWidth) - 1n))).toLocaleString()} 〜 ${((1n << (BigInt(bitWidth) - 1n)) - 1n).toLocaleString()}`
              : `0 〜 ${((1n << BigInt(bitWidth)) - 1n).toLocaleString()}`}
          </div>
        </div>

        {/* 四則演算セクション */}
        <div className="border-t border-gray-600 pt-8">
          <h2 className="text-xl font-bold mb-1">四則演算 / ビット演算</h2>
          <p className="text-gray-400 text-xs mb-5">各オペランドの基数を選んで計算。結果は丸めずにそのまま表示します。</p>

          {/* オペランド A */}
          <div className="flex gap-2 items-end mb-3">
            <div className="flex-1">
              <label className="block text-sm font-semibold mb-1 text-gray-300">A</label>
              <input
                type="text"
                value={calcA}
                onChange={(e) => handleCalcAChange(e.target.value)}
                className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono"
                placeholder={calcInputPlaceholder(calcBaseA)}
                spellCheck={false}
              />
            </div>
            <div>
              <label className="block text-xs font-semibold mb-1 text-gray-400">基数</label>
              <select
                value={calcBaseA}
                onChange={(e) => setCalcBaseA(Number(e.target.value) as CalcBase)}
                className="p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 text-sm"
              >
                <option value={2}>2進数</option>
                <option value={10}>10進数</option>
                <option value={16}>16進数</option>
              </select>
            </div>
          </div>

          {/* 演算子 */}
          <div className="mb-3">
            <label className="block text-sm font-semibold mb-2 text-gray-300">演算子</label>
            <div className="flex flex-wrap gap-2">
              {OPERATORS.map((op) => (
                <button
                  key={op.value}
                  type="button"
                  title={op.title}
                  onClick={() => setCalcOp(op.value)}
                  className={`px-3 py-2 rounded font-mono font-bold text-sm transition min-w-[2.5rem] ${calcOp === op.value
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-700 hover:bg-gray-600 text-gray-200'
                    }`}
                >
                  {op.label}
                </button>
              ))}
            </div>
            <p className="text-xs text-gray-500 mt-1">
              {OPERATORS.find((o) => o.value === calcOp)?.title}
            </p>
          </div>

          {/* オペランド B */}
          <div className="flex gap-2 items-end mb-5">
            <div className="flex-1">
              <label className="block text-sm font-semibold mb-1 text-gray-300">B</label>
              <input
                type="text"
                value={calcB}
                onChange={(e) => handleCalcBChange(e.target.value)}
                className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono"
                placeholder={calcInputPlaceholder(calcBaseB)}
                spellCheck={false}
              />
            </div>
            <div>
              <label className="block text-xs font-semibold mb-1 text-gray-400">基数</label>
              <select
                value={calcBaseB}
                onChange={(e) => setCalcBaseB(Number(e.target.value) as CalcBase)}
                className="p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 text-sm"
              >
                <option value={2}>2進数</option>
                <option value={10}>10進数</option>
                <option value={16}>16進数</option>
              </select>
            </div>
          </div>

          {/* 計算エラー */}
          {calcError && (
            <div className="mb-4 p-3 rounded bg-red-900 border border-red-700 text-red-300 text-sm">
              {calcError}
            </div>
          )}

          {/* 計算結果 */}
          {calcResult && (
            <div className="bg-gray-700 rounded-lg p-4 space-y-3">
              <div className="text-xs font-semibold text-gray-400 mb-1">計算結果</div>
              <div className="flex items-start gap-3">
                <span className="text-xs text-gray-400 w-20 shrink-0 pt-1">2進数</span>
                <span className="font-mono text-green-400 break-all leading-relaxed">{calcResult.grouped}</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="text-xs text-gray-400 w-20 shrink-0">10進数</span>
                <span className="font-mono text-yellow-300 break-all">{calcResult.dec}</span>
              </div>
              <div className="flex items-center gap-3">
                <span className="text-xs text-gray-400 w-20 shrink-0">16進数</span>
                <span className="font-mono text-cyan-300 break-all">0x{calcResult.hex}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
