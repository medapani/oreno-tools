import React, { useState } from 'react';
import { useSessionState } from '../hooks/useSessionState';
import * as yaml from 'js-yaml';

export default function JSONFormatter() {
  const [input, setInput] = useSessionState('json.input', '');
  const [output, setOutput] = useSessionState('json.output', '');
  const [error, setError] = useState('');
  const [indentSize, setIndentSize] = useSessionState('json.indentSize', '2');

  // 入力がJSONかYAMLかを判定する
  const detectFormat = (text: string): 'json' | 'yaml' | 'unknown' => {
    if (!text.trim()) return 'unknown';

    // まずJSONとしてパースを試みる
    try {
      JSON.parse(text);
      return 'json';
    } catch {
      // JSONでない場合、YAMLとして試みる
      try {
        yaml.load(text);
        return 'yaml';
      } catch {
        return 'unknown';
      }
    }
  };

  const handleFormat = () => {
    const format = detectFormat(input);

    if (format === 'unknown') {
      setOutput('');
      setError('入力がJSON/YAMLどちらとしても解析できませんでした');
      return;
    }

    try {
      let parsed: any;
      if (format === 'json') {
        parsed = JSON.parse(input);
        const formatted = JSON.stringify(parsed, null, parseInt(indentSize));
        setOutput(formatted);
      } else {
        // YAML
        parsed = yaml.load(input);
        const formatted = yaml.dump(parsed, {
          indent: parseInt(indentSize),
          lineWidth: -1,
          noRefs: true,
        });
        setOutput(formatted);
      }
      setError('');
    } catch (e) {
      setOutput('');
      if (e instanceof Error) {
        setError(`整形に失敗しました: ${e.message}`);
      } else {
        setError(`整形に失敗しました: ${String(e)}`);
      }
    }
  };

  const handleMinify = () => {
    const format = detectFormat(input);

    if (format === 'unknown') {
      setOutput('');
      setError('入力がJSON/YAMLどちらとしても解析できませんでした');
      return;
    }

    try {
      let parsed: any;
      if (format === 'json') {
        parsed = JSON.parse(input);
        const minified = JSON.stringify(parsed);
        setOutput(minified);
      } else {
        // YAMLの場合は最小限の設定で出力
        parsed = yaml.load(input);
        const minified = yaml.dump(parsed, {
          indent: 2,
          lineWidth: 120,
          noRefs: true,
          flowLevel: 0, // インライン形式を使用
        });
        setOutput(minified);
      }
      setError('');
    } catch (e) {
      setOutput('');
      if (e instanceof Error) {
        setError(`圧縮に失敗しました: ${e.message}`);
      } else {
        setError(`圧縮に失敗しました: ${String(e)}`);
      }
    }
  };

  const handleValidate = () => {
    const format = detectFormat(input);

    if (format === 'unknown') {
      setOutput('');
      setError('✗ 入力がJSON/YAMLどちらとしても解析できませんでした');
      return;
    }

    try {
      if (format === 'json') {
        JSON.parse(input);
        setError('');
        setOutput('✓ 有効なJSONです');
      } else {
        yaml.load(input);
        setError('');
        setOutput('✓ 有効なYAMLです');
      }
    } catch (e) {
      setOutput('');
      if (e instanceof Error) {
        setError(`✗ 無効なデータ: ${e.message}`);
      } else {
        setError(`✗ 無効なデータ: ${String(e)}`);
      }
    }
  };

  const handleClear = () => {
    setInput('');
    setOutput('');
    setError('');
  };

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(output);
      // コピー成功の視覚的フィードバック（オプション）
    } catch (e) {
      setError('クリップボードへのコピーに失敗しました');
    }
  };

  const handleJsonToYaml = () => {
    try {
      const parsed = JSON.parse(input);
      const yamlStr = yaml.dump(parsed, {
        indent: parseInt(indentSize),
        lineWidth: -1, // 行の折り返しを無効化
        noRefs: true,
      });
      setOutput(yamlStr);
      setError('');
    } catch (e) {
      setOutput('');
      if (e instanceof Error) {
        setError(`YAML変換に失敗しました: ${e.message}`);
      } else {
        setError(`YAML変換に失敗しました: ${String(e)}`);
      }
    }
  };

  const handleYamlToJson = () => {
    try {
      const parsed = yaml.load(input);
      const jsonStr = JSON.stringify(parsed, null, parseInt(indentSize));
      setOutput(jsonStr);
      setError('');
    } catch (e) {
      setOutput('');
      if (e instanceof Error) {
        setError(`JSON変換に失敗しました: ${e.message}`);
      } else {
        setError(`JSON変換に失敗しました: ${String(e)}`);
      }
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-6xl">
        <h1 className="text-3xl font-bold mb-6 text-center">JSON/YAML 変換</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* 入力セクション */}
          <div>
            <div className="flex justify-between items-center mb-2">
              <label htmlFor="json-input" className="block text-sm font-semibold">
                JSON/YAML 入力
              </label>
              <div className="flex items-center gap-2">
                <label htmlFor="indent-size" className="text-sm font-semibold">
                  インデント:
                </label>
                <select
                  id="indent-size"
                  value={indentSize}
                  onChange={(e) => setIndentSize(e.target.value)}
                  className="p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 text-sm"
                >
                  <option value="2">2スペース</option>
                  <option value="4">4スペース</option>
                  <option value="8">8スペース</option>
                  <option value="1">タブ</option>
                </select>
              </div>
            </div>
            <textarea
              id="json-input"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              className="w-full h-96 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-sm"
              placeholder='JSON: {"key": "value", "array": [1, 2, 3]}&#10;&#10;YAML:&#10;key: value&#10;array:&#10;  - 1&#10;  - 2&#10;  - 3'
            />
          </div>

          {/* 出力セクション */}
          <div>
            <div className="flex justify-between items-center mb-2">
              <label htmlFor="json-output" className="block text-sm font-semibold">
                JSON/YAML 出力
              </label>
              {output && (
                <button
                  onClick={handleCopy}
                  className="px-3 py-1 rounded bg-gray-600 hover:bg-gray-500 transition text-sm font-semibold"
                  type="button"
                  title="クリップボードにコピー"
                >
                  コピー
                </button>
              )}
            </div>
            <textarea
              id="json-output"
              value={output}
              readOnly
              className="w-full h-96 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none font-mono text-sm"
              placeholder="整形されたJSON/YAMLがここに表示されます"
            />
          </div>
        </div>

        {/* エラー表示 */}
        {error && (
          <div className="mt-6 p-3 bg-red-900 border border-red-600 rounded text-red-100">
            {error}
          </div>
        )}

        {/* アクションボタン */}
        <div className="mt-6 flex flex-wrap gap-3 justify-center">
          <button
            onClick={handleFormat}
            className="px-6 py-2 rounded bg-blue-600 hover:bg-blue-500 transition font-semibold"
            type="button"
          >
            整形
          </button>
          <button
            onClick={handleMinify}
            className="px-6 py-2 rounded bg-green-600 hover:bg-green-500 transition font-semibold"
            type="button"
          >
            圧縮
          </button>
          <button
            onClick={handleValidate}
            className="px-6 py-2 rounded bg-purple-600 hover:bg-purple-500 transition font-semibold"
            type="button"
          >
            検証
          </button>
          <button
            onClick={handleJsonToYaml}
            className="px-6 py-2 rounded bg-yellow-600 hover:bg-yellow-500 transition font-semibold"
            type="button"
          >
            JSON → YAML
          </button>
          <button
            onClick={handleYamlToJson}
            className="px-6 py-2 rounded bg-orange-600 hover:bg-orange-500 transition font-semibold"
            type="button"
          >
            YAML → JSON
          </button>
          <button
            onClick={handleClear}
            className="px-6 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
            type="button"
          >
            クリア
          </button>
        </div>
      </div>
    </div>
  );
}
