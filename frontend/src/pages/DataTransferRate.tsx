// データ転送速度（B/s, MB/s, Mbps等）に変換するアプリ
import React, { useState, useEffect } from 'react';
import { ConvertDataTransferRate } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

export default function DataTransferRate() {
  const [input, setInput] = useSessionState('dataTransferRate.input', '');
  const [result, setResult] = useState('');
  const [unit, setUnit] = useSessionState('dataTransferRate.unit', 'Mbps');

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
          const result = await ConvertDataTransferRate(value, unit);

          // 結果を整形して表示
          const unitOrder = [
            { key: 'B/s', label: 'B/s' },
            { key: 'KB/s', label: 'KB/s' },
            { key: 'MB/s', label: 'MB/s' },
            { key: 'GB/s', label: 'GB/s' },
            { key: 'TB/s', label: 'TB/s' },
            { key: 'KiB/s', label: 'KiB/s' },
            { key: 'MiB/s', label: 'MiB/s' },
            { key: 'GiB/s', label: 'GiB/s' },
            { key: 'TiB/s', label: 'TiB/s' },
            { key: 'bit/s', label: 'bps' },
            { key: 'Kbit/s', label: 'Kbps' },
            { key: 'Mbit/s', label: 'Mbps' },
            { key: 'Gbit/s', label: 'Gbps' },
            { key: 'Tbit/s', label: 'Tbps' },
            { key: 'Kibit/s', label: 'Kibps' },
            { key: 'Mibit/s', label: 'Mibps' },
            { key: 'Gibit/s', label: 'Gibps' },
            { key: 'Tibit/s', label: 'Tibps' }
          ];

          const formatted = unitOrder
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
          setResult(formatted);
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
        <h1 className="text-3xl font-bold mb-6 text-center">データ転送速度変換</h1>

        <input
          type="text"
          placeholder="転送速度を入力してください (例: 100)"
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
          <option value="B/s">B/s</option>
          <option value="KB/s">KB/s</option>
          <option value="MB/s">MB/s</option>
          <option value="GB/s">GB/s</option>
          <option value="TB/s">TB/s</option>
          <option value="KiB/s">KiB/s</option>
          <option value="MiB/s">MiB/s</option>
          <option value="GiB/s">GiB/s</option>
          <option value="TiB/s">TiB/s</option>
          <option value="bps">bps</option>
          <option value="Kbps">Kbps</option>
          <option value="Mbps">Mbps</option>
          <option value="Gbps">Gbps</option>
          <option value="Tbps">Tbps</option>
          <option value="Kibps">Kibps</option>
          <option value="Mibps">Mibps</option>
          <option value="Gibps">Gibps</option>
          <option value="Tibps">Tibps</option>
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
