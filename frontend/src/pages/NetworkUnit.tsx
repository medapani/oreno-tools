// 通信単位（bits, Kbit, Mbit等）に変換するアプリ
import React, { useState, useEffect } from 'react';
import { ConvertUnit } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

export default function NetworkUnit() {
  const [input, setInput] = useSessionState('networkUnit.input', '');
  const [result, setResult] = useState('');
  const [unit, setUnit] = useSessionState('networkUnit.unit', 'Mbps');

  useEffect(() => {
    const convertValue = async () => {
      if (!input) {
        setResult('');
        return;
      }

      const match = input.match(/^(\d+\.?\d*)$/);
      if (match) {
        const value = parseFloat(match[1]);

        try {
          const backendUnitMap: Record<string, string> = {
            Kbps: 'Kbits',
            Mbps: 'Mbits',
            Gbps: 'Gbits'
          };
          const backendUnit = backendUnitMap[unit] || 'Mbits';
          const result = await ConvertUnit(value, backendUnit);

          // 結果を整形して表示
          const bitUnitOrder = [
            { key: 'Kbits', label: 'Kbps' },
            { key: 'Mbits', label: 'Mbps' },
            { key: 'Gbits', label: 'Gbps' }
          ];
          const byteUnitOrder = [
            { key: 'B', label: 'B/s' },
            { key: 'KB', label: 'KB/s' },
            { key: 'MB', label: 'MB/s' },
            { key: 'GB', label: 'GB/s' }
          ];

          const bitFormatted = bitUnitOrder
            .map(({ key, label }) => {
              const val = result[key as keyof typeof result];
              const numVal = Number(val);
              const formattedNum = numVal.toLocaleString('ja-JP', {
                minimumFractionDigits: 0,
                maximumFractionDigits: 3
              });
              return `${label}: ${formattedNum}`;
            })
            .join('\n');

          const byteFormatted = byteUnitOrder
            .map(({ key, label }) => {
              const val = result[key as keyof typeof result];
              const numVal = Number(val);
              const formattedNum = numVal.toLocaleString('ja-JP', {
                minimumFractionDigits: 0,
                maximumFractionDigits: 3
              });
              return `${label}: ${formattedNum}`;
            })
            .join('\n');

          setResult(`${bitFormatted}\n\n${byteFormatted}`);
        } catch (error) {
          setResult(`エラー: ${error}`);
        }
      } else {
        setResult('無効な入力です（数値のみ入力してください）');
      }
    };

    convertValue();
  }, [input, unit]);

  return (
    <div className="h-screen flex items-center justify-center bg-gray-900 text-white">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-96">
        <h1 className="text-3xl font-bold mb-6 text-center">通信速度変換</h1>

        <input
          type="text"
          placeholder="通信速度を入力してください (例: 100)"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          className="w-full p-3 mb-4 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
        />

        <select
          name="unit"
          id="unit"
          value={unit}
          onChange={(e) => setUnit(e.target.value)}
          className="w-full p-3 mb-4 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
        >
          <option value="Kbps">Kbps</option>
          <option value="Mbps">Mbps</option>
          <option value="Gbps">Gbps</option>
        </select>

        {result && (
          <div className="mt-4 p-3 bg-gray-700 rounded whitespace-pre-line">
            結果: <br />{result}
          </div>
        )}
      </div>
    </div>
  );
}
