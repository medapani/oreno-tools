import React, { useState, useEffect, useRef } from 'react';
import jsQR from 'jsqr';
import { GenerateQRCode, SaveQRCodePNGAs } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

const SIZE_OPTIONS = [128, 256, 512] as const;
const LEVEL_OPTIONS = [
  { value: 'L', label: 'L（低 7%）' },
  { value: 'M', label: 'M（中 15%）' },
  { value: 'Q', label: 'Q（高 25%）' },
  { value: 'H', label: 'H（最高 30%）' },
] as const;

export default function QRCodeTool() {
  const [input, setInput] = useSessionState('qrcode.input', '');
  const [size, setSize] = useSessionState<number>('qrcode.size', 256);
  const [level, setLevel] = useSessionState('qrcode.level', 'M');
  const [qrImage, setQrImage] = useState('');
  const [decodedText, setDecodedText] = useState('');
  const [decodeError, setDecodeError] = useState('');
  const [decodeImagePreview, setDecodeImagePreview] = useState('');
  const [isDecoding, setIsDecoding] = useState(false);
  const [isDragOver, setIsDragOver] = useState(false);
  const [copyMessage, setCopyMessage] = useState('');
  const [error, setError] = useState('');
  const requestIdRef = useRef(0);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const decodeFromDataURL = async (dataURL: string) => {
    setIsDecoding(true);
    setDecodeError('');

    try {
      const image = await new Promise<HTMLImageElement>((resolve, reject) => {
        const img = new Image();
        img.onload = () => resolve(img);
        img.onerror = () => reject(new Error('画像の読み込みに失敗しました'));
        img.src = dataURL;
      });

      const canvas = document.createElement('canvas');
      canvas.width = image.width;
      canvas.height = image.height;
      const context = canvas.getContext('2d');
      if (!context) {
        throw new Error('画像解析の初期化に失敗しました');
      }

      context.drawImage(image, 0, 0);
      const imageData = context.getImageData(0, 0, canvas.width, canvas.height);
      const result = jsQR(imageData.data, imageData.width, imageData.height, {
        inversionAttempts: 'attemptBoth',
      });

      if (!result) {
        throw new Error('QRコードを検出できませんでした');
      }

      setDecodedText(result.data);
      setDecodeImagePreview(dataURL);
      setDecodeError('');
    } catch (e) {
      setDecodedText('');
      setDecodeImagePreview(dataURL);
      setDecodeError(String(e));
    } finally {
      setIsDecoding(false);
    }
  };

  const decodeFromFile = async (file: File) => {
    if (!file.type.startsWith('image/')) {
      setDecodedText('');
      setDecodeImagePreview('');
      setDecodeError('画像ファイルを選択してください');
      return;
    }

    const dataURL = await new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        if (typeof reader.result === 'string') {
          resolve(reader.result);
          return;
        }
        reject(new Error('ファイルの読み取りに失敗しました'));
      };
      reader.onerror = () => reject(new Error('ファイルの読み取りに失敗しました'));
      reader.readAsDataURL(file);
    });

    await decodeFromDataURL(dataURL);
  };

  useEffect(() => {
    const onPaste = (event: ClipboardEvent) => {
      const items = event.clipboardData?.items;
      if (!items) {
        return;
      }

      for (const item of items) {
        if (item.type.startsWith('image/')) {
          const file = item.getAsFile();
          if (file) {
            event.preventDefault();
            void decodeFromFile(file);
          }
          return;
        }
      }
    };

    document.addEventListener('paste', onPaste);
    return () => document.removeEventListener('paste', onPaste);
  }, []);

  useEffect(() => {
    if (!input.trim()) {
      setQrImage('');
      setError('');
      return;
    }

    const currentId = ++requestIdRef.current;
    const timer = setTimeout(async () => {
      try {
        const result = await GenerateQRCode(input, size, level);
        if (currentId === requestIdRef.current) {
          setQrImage(result);
          setError('');
        }
      } catch (e) {
        if (currentId === requestIdRef.current) {
          setQrImage('');
          setError(String(e));
        }
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [input, size, level]);

  useEffect(() => {
    if (!copyMessage) {
      return;
    }

    const timer = window.setTimeout(() => {
      setCopyMessage('');
    }, 3000);

    return () => window.clearTimeout(timer);
  }, [copyMessage]);

  const handleDownloadAs = async () => {
    if (!qrImage) return;

    try {
      const savedPath = await SaveQRCodePNGAs(qrImage, 'qrcode.png');
      if (!savedPath) {
        return;
      }
      setError('');
      alert(`保存しました: ${savedPath}`);
    } catch (e) {
      setError(`PNG 保存に失敗しました: ${String(e)}`);
    }
  };

  const handleClear = () => {
    setInput('');
    setQrImage('');
    setError('');
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }
    void decodeFromFile(file);
    event.target.value = '';
  };

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    setIsDragOver(false);
    const file = event.dataTransfer.files?.[0];
    if (!file) {
      return;
    }
    void decodeFromFile(file);
  };

  const handleCopyDecodedText = async () => {
    if (!decodedText.trim()) {
      setCopyMessage('コピーするテキストがありません');
      return;
    }

    try {
      await navigator.clipboard.writeText(decodedText);
      setCopyMessage('解析結果をコピーしました');
    } catch (_e) {
      setCopyMessage('コピーに失敗しました');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-2xl">
        <h1 className="text-3xl font-bold mb-6 text-center">QR コード生成</h1>

        <div className="mb-4">
          <label htmlFor="qr-input" className="block text-sm font-semibold mb-2">
            テキスト / URL
          </label>
          <textarea
            id="qr-input"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            className="w-full h-28 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-sm"
            placeholder="QR コードにするテキストまたは URL を入力"
          />
        </div>

        <div className="flex flex-wrap gap-6 mb-6">
          <div>
            <label className="block text-sm font-semibold mb-2">サイズ (px)</label>
            <select
              value={size}
              onChange={(e) => setSize(Number(e.target.value))}
              className="p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            >
              {SIZE_OPTIONS.map((s) => (
                <option key={s} value={s}>{s} × {s}</option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-semibold mb-2">誤り訂正レベル</label>
            <select
              value={level}
              onChange={(e) => setLevel(e.target.value)}
              className="p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            >
              {LEVEL_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </select>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-900 border border-red-600 rounded text-red-100 text-sm">
            {error}
          </div>
        )}

        {qrImage && (
          <div className="flex flex-col items-center gap-4">
            <img
              src={`data:image/png;base64,${qrImage}`}
              alt="Generated QR Code"
              className="rounded border border-gray-600 bg-white p-2"
              style={{ width: size, height: size }}
            />
            <div className="flex gap-3">
              <button
                onClick={() => void handleDownloadAs()}
                className="px-5 py-2 rounded bg-cyan-600 hover:bg-cyan-500 transition font-semibold"
                type="button"
              >
                保存先を選んで保存
              </button>
              <button
                onClick={handleClear}
                className="px-5 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
                type="button"
              >
                クリア
              </button>
            </div>
          </div>
        )}

        {!qrImage && !error && (
          <div className="flex justify-center mt-4">
            <button
              onClick={handleClear}
              className="px-5 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
              type="button"
            >
              クリア
            </button>
          </div>
        )}

        <hr className="my-8 border-gray-700" />

        <h2 className="text-2xl font-bold mb-4 text-center">QR コード解析</h2>

        <div
          onDragOver={(event) => {
            event.preventDefault();
            setIsDragOver(true);
          }}
          onDragLeave={() => setIsDragOver(false)}
          onDrop={handleDrop}
          className={`rounded-lg border-2 border-dashed p-6 text-center transition ${isDragOver ? 'border-blue-400 bg-blue-900/20' : 'border-gray-600 bg-gray-700/30'
            }`}
        >
          <p className="text-sm text-gray-200 mb-3">
            画像をドラッグ&ドロップ、または下のボタンで選択してください
          </p>
          <p className="text-xs text-gray-400 mb-4">この画面上で画像をペースト (Cmd+V / Ctrl+V) しても解析できます</p>

          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={handleFileChange}
          />
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="px-4 py-2 rounded bg-indigo-600 hover:bg-indigo-500 transition font-semibold"
          >
            画像ファイルを選択
          </button>
        </div>

        {isDecoding && (
          <div className="mt-4 p-3 rounded bg-gray-700 text-sm text-gray-100">解析中...</div>
        )}

        {decodeError && (
          <div className="mt-4 p-3 rounded bg-red-900 border border-red-600 text-sm text-red-100">{decodeError}</div>
        )}

        {decodeImagePreview && (
          <div className="mt-6 flex flex-col gap-4">
            <img
              src={decodeImagePreview}
              alt="Uploaded for QR decode"
              className="max-h-80 object-contain rounded border border-gray-600 bg-white p-2"
            />

            <div>
              <label htmlFor="decoded-text" className="block text-sm font-semibold mb-2">
                解析結果テキスト
              </label>
              <textarea
                id="decoded-text"
                value={decodedText}
                readOnly
                className="w-full h-28 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-sm"
                placeholder="解析結果がここに表示されます"
              />
              <div className="mt-3 flex items-center gap-3">
                <button
                  type="button"
                  onClick={() => void handleCopyDecodedText()}
                  className="px-4 py-2 rounded bg-emerald-600 hover:bg-emerald-500 transition font-semibold"
                >
                  解析結果をコピー
                </button>
                {copyMessage && <p className="text-sm text-gray-300">{copyMessage}</p>}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
