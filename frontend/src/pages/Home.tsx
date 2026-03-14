// B -> MB, MiB -> MiB, GB -> GB, GiB -> GiBなどに変換するアプリ
import React, { useState, useEffect } from 'react';
import { ConvertUnit } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

export default function Home() {
  const [input, setInput] = useSessionState('home.input', '');
  const [result, setResult] = useState('');
  const [unit, setUnit] = useSessionState('home.unit', 'B');

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
          const result = await ConvertUnit(value, unit);

          // 結果を整形して表示
          const unitOrder = ['B', 'KB', 'MB', 'GB', 'TB', 'KiB', 'MiB', 'GiB', 'TiB'];
          const formatted = unitOrder
            .map(key => {
              const val = result[key as keyof typeof result];
              const numVal = Number(val);
              const formattedNum = numVal.toLocaleString('ja-JP', {
                minimumFractionDigits: 0,
                maximumFractionDigits: 3
              });
              return `${key}: ${formattedNum}`;
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
        <h1 className="text-3xl font-bold mb-6 text-center">バイト変換</h1>

        <input
          type="text"
          placeholder="容量を入力してください (例: 10)"
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
          <option value="B">B</option>
          <option value="KB">KB</option>
          <option value="MB">MB</option>
          <option value="GB">GB</option>
          <option value="TB">TB</option>
          <option value="KiB">KiB</option>
          <option value="MiB">MiB</option>
          <option value="GiB">GiB</option>
          <option value="TiB">TiB</option>
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