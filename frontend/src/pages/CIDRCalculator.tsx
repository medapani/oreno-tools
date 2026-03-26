import React, { useState, useEffect } from 'react';
import { CalculateCIDR } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

interface CIDRResult {
  networkAddress: string;
  broadcastAddress: string;
  subnetMask: string;
  wildcardMask: string;
  firstHostAddress: string;
  lastHostAddress: string;
  totalHosts: number;
  usableHosts: number;
  cidr: string;
  binarySubnetMask: string;
  ipClass: string;
  ipType: string;
  inputIp: string;
  inputWasHost: boolean;
}

function buildInputNotice(result: CIDRResult): string {
  return `入力IP ${result.inputIp} はネットワークアドレスではなくホストアドレスです。単一ホストを表す場合は ${result.inputIp}/32 を使用してください。`;
}

export default function CIDRCalculator() {
  const [input, setInput] = useSessionState('cidrCalculator.input', '192.168.1.0/24');
  const [result, setResult] = useState<CIDRResult | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    const calculate = async () => {
      if (!input) {
        setResult(null);
        setError('');
        return;
      }

      try {
        const res = await CalculateCIDR(input);
        setResult(res);
        setError('');
      } catch (err) {
        setError(err instanceof Error ? err.message : String(err));
        setResult(null);
      }
    };

    calculate();
  }, [input]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-8 rounded-xl shadow-2xl w-full max-w-2xl">
        <h1 className="text-3xl font-bold mb-6 text-center">CIDR 計算</h1>

        <div className="mb-6">
          <label className="block text-sm font-medium mb-2">
            CIDR表記を入力してください
          </label>
          <input
            type="text"
            placeholder="例: 192.168.1.0/24"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono"
          />
          <p className="text-xs text-gray-400 mt-2">
            例: 192.168.1.0/24, 10.0.0.0/8, 172.16.0.0/12
          </p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-900 border border-red-700 rounded">
            <p className="text-red-200">{error}</p>
          </div>
        )}

        {result && (
          <div className="space-y-4">
            {result.inputWasHost && (
              <div className="p-4 bg-amber-900/60 border border-amber-600 rounded">
                <div className="text-amber-300 font-semibold mb-1">入力値に関する注意</div>
                <p className="text-amber-100 text-sm leading-relaxed">{buildInputNotice(result)}</p>
              </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">ネットワークアドレス</div>
                <div className="font-mono text-lg">{result.networkAddress}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">ブロードキャストアドレス</div>
                <div className="font-mono text-lg">{result.broadcastAddress}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">サブネットマスク</div>
                <div className="font-mono text-lg">{result.subnetMask}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">ワイルドカードマスク</div>
                <div className="font-mono text-lg">{result.wildcardMask}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">最初のホストアドレス</div>
                <div className="font-mono text-lg">{result.firstHostAddress}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">最後のホストアドレス</div>
                <div className="font-mono text-lg">{result.lastHostAddress}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">総ホスト数</div>
                <div className="font-mono text-lg">{result.totalHosts.toLocaleString('ja-JP')}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">使用可能ホスト数</div>
                <div className="font-mono text-lg">{result.usableHosts.toLocaleString('ja-JP')}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">IPクラス</div>
                <div className="font-mono text-lg">{result.ipClass}</div>
              </div>

              <div className="bg-gray-700 p-4 rounded">
                <div className="text-sm text-gray-400 mb-1">IPタイプ</div>
                <div className="font-mono text-lg">{result.ipType}</div>
              </div>
            </div>

            <div className="bg-gray-700 p-4 rounded">
              <div className="text-sm text-gray-400 mb-1">バイナリ サブネットマスク</div>
              <div className="font-mono text-sm break-all">{result.binarySubnetMask}</div>
            </div>

            <div className="bg-gray-700 p-4 rounded">
              <div className="text-sm text-gray-400 mb-2">サマリー</div>
              <div className="text-sm space-y-1">
                <div>入力IP: <span className="font-mono">{result.inputIp}</span></div>
                <div>CIDR: <span className="font-mono">{result.cidr}</span></div>
                <div>ネットワーク範囲: <span className="font-mono">
                  {result.firstHostAddress} - {result.lastHostAddress}
                </span></div>
                <div>使用可能IP数: <span className="font-mono">
                  {result.usableHosts.toLocaleString('ja-JP')} 個
                </span></div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
