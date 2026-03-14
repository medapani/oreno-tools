import React, { useState } from 'react';
import { AddCertificatesToCRL, SaveTextFile } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

function extractPemCertificates(input: string): string[] {
  const matches = input.match(/-----BEGIN CERTIFICATE-----[\s\S]*?-----END CERTIFICATE-----/g);
  if (matches && matches.length > 0) {
    return matches.map((m) => m.trim()).filter((m) => m.length > 0);
  }

  const trimmed = input.trim();
  return trimmed ? [trimmed] : [];
}

export default function CRLTool() {
  const [caCertificatePem, setCaCertificatePem] = useSessionState('crl.caCertificatePem', '');
  const [caPrivateKeyPem, setCaPrivateKeyPem] = useSessionState('crl.caPrivateKeyPem', '');
  const [existingCrlPem, setExistingCrlPem] = useSessionState('crl.existingCrlPem', '');
  const [revokedCertsInput, setRevokedCertsInput] = useSessionState('crl.revokedCertsInput', '');
  const [nextUpdateDays, setNextUpdateDays] = useSessionState('crl.nextUpdateDays', '30');

  const [crlPem, setCrlPem] = useSessionState('crl.crlPem', '');
  const [addedCount, setAddedCount] = useSessionState('crl.addedCount', 0);
  const [totalRevokedCount, setTotalRevokedCount] = useSessionState('crl.totalRevokedCount', 0);
  const [revokedSerialNumbers, setRevokedSerialNumbers] = useSessionState<string[]>('crl.revokedSerialNumbers', []);
  const [savedPaths, setSavedPaths] = useState<string[]>([]);

  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleGenerate = async () => {
    setError('');
    setCrlPem('');
    setAddedCount(0);
    setTotalRevokedCount(0);
    setRevokedSerialNumbers([]);
    setSavedPaths([]);

    const days = Number(nextUpdateDays);
    if (!Number.isFinite(days) || days <= 0 || !Number.isInteger(days)) {
      setError('CRLの次回更新日数は1以上の整数で入力してください');
      return;
    }

    const revokedCerts = extractPemCertificates(revokedCertsInput);
    if (revokedCerts.length === 0) {
      setError('失効対象の証明書(PEM)を1つ以上入力してください');
      return;
    }

    try {
      setIsLoading(true);
      const result = await AddCertificatesToCRL(
        caCertificatePem,
        caPrivateKeyPem,
        existingCrlPem,
        revokedCerts,
        days
      );

      setCrlPem(result.crlPem);
      setAddedCount(result.addedCount);
      setTotalRevokedCount(result.totalRevokedCount);
      setRevokedSerialNumbers(result.revokedSerialNumbers || []);
    } catch (e) {
      setError(`CRLの更新に失敗しました: ${String(e)}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = (text: string, label: string) => {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        alert(`${label}をコピーしました`);
      })
      .catch(() => {
        alert(`${label}のコピーに失敗しました`);
      });
  };

  const handleDownload = async (content: string, filename: string, notify = true) => {
    try {
      const savedPath = await SaveTextFile(content, filename);
      setSavedPaths((prev) => {
        const next = [savedPath, ...prev.filter((p) => p !== savedPath)];
        return next.slice(0, 30);
      });
      if (notify) {
        alert(`保存しました: ${savedPath}`);
      }
      return savedPath;
    } catch (e) {
      setError(`ファイル保存に失敗しました: ${String(e)}`);
      return '';
    }
  };

  const handleClear = () => {
    setCaCertificatePem('');
    setCaPrivateKeyPem('');
    setExistingCrlPem('');
    setRevokedCertsInput('');
    setNextUpdateDays('30');
    setCrlPem('');
    setAddedCount(0);
    setTotalRevokedCount(0);
    setRevokedSerialNumbers([]);
    setSavedPaths([]);
    setError('');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-6xl">
        <h1 className="text-3xl font-bold mb-6 text-center">CRL 更新</h1>
        <p className="text-gray-400 text-center mb-6">
          既存CRLに失効対象証明書を追加して、新しいCRLを再発行します
        </p>

        <div className="space-y-4 mb-6">
          <div>
            <label htmlFor="ca-cert-pem" className="block text-sm font-semibold mb-2">CA証明書 (PEM)</label>
            <textarea
              id="ca-cert-pem"
              value={caCertificatePem}
              onChange={(e) => setCaCertificatePem(e.target.value)}
              className="w-full h-36 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs focus:outline-none focus:border-blue-500"
              placeholder="-----BEGIN CERTIFICATE-----"
            />
          </div>

          <div>
            <label htmlFor="ca-key-pem" className="block text-sm font-semibold mb-2">CA秘密鍵 (PEM)</label>
            <textarea
              id="ca-key-pem"
              value={caPrivateKeyPem}
              onChange={(e) => setCaPrivateKeyPem(e.target.value)}
              className="w-full h-36 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs focus:outline-none focus:border-blue-500"
              placeholder="-----BEGIN PRIVATE KEY-----"
            />
          </div>

          <div>
            <label htmlFor="existing-crl-pem" className="block text-sm font-semibold mb-2">既存CRL (PEM, 任意)</label>
            <textarea
              id="existing-crl-pem"
              value={existingCrlPem}
              onChange={(e) => setExistingCrlPem(e.target.value)}
              className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs focus:outline-none focus:border-blue-500"
              placeholder="-----BEGIN X509 CRL-----"
            />
          </div>

          <div>
            <label htmlFor="revoked-certs" className="block text-sm font-semibold mb-2">
              失効対象証明書 (PEM, 複数可)
            </label>
            <textarea
              id="revoked-certs"
              value={revokedCertsInput}
              onChange={(e) => setRevokedCertsInput(e.target.value)}
              className="w-full h-48 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs focus:outline-none focus:border-blue-500"
              placeholder={"-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"}
            />
          </div>

          <div>
            <label htmlFor="next-update-days" className="block text-sm font-semibold mb-2">CRL次回更新日数</label>
            <input
              id="next-update-days"
              type="number"
              min={1}
              value={nextUpdateDays}
              onChange={(e) => setNextUpdateDays(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            />
          </div>
        </div>

        <div className="flex flex-col gap-3 mb-6">
          <button
            onClick={handleGenerate}
            disabled={isLoading || !caCertificatePem.trim() || !caPrivateKeyPem.trim() || !revokedCertsInput.trim()}
            className="w-full px-4 py-3 rounded bg-blue-600 hover:bg-blue-500 disabled:bg-gray-500 disabled:cursor-not-allowed transition font-semibold text-lg"
            type="button"
          >
            {isLoading ? 'CRL更新中...' : 'CRLに証明書を追加'}
          </button>
          {crlPem && (
            <button
              onClick={() => handleDownload(crlPem, 'ca-updated.crl.pem')}
              className="w-full px-4 py-3 rounded bg-green-600 hover:bg-green-500 transition font-semibold text-lg"
              type="button"
            >
              CRLを保存
            </button>
          )}
          <button
            onClick={handleClear}
            className="w-full px-4 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
            type="button"
          >
            クリア
          </button>
        </div>

        {savedPaths.length > 0 && (
          <div className="mb-6 p-4 rounded border border-emerald-500 bg-emerald-900/30">
            <div className="font-semibold text-emerald-300 mb-2">保存先</div>
            <div className="text-xs text-emerald-100 break-all max-h-40 overflow-auto space-y-1">
              {savedPaths.map((path, index) => (
                <div key={`${path}-${index}`}>{path}</div>
              ))}
            </div>
          </div>
        )}

        {error && (
          <div className="mb-6 p-4 bg-red-600 rounded text-white">
            <strong>エラー:</strong> {error}
          </div>
        )}

        {crlPem && (
          <div className="space-y-6">
            <div className="border border-gray-600 rounded p-4">
              <h2 className="text-xl font-bold mb-3 text-cyan-400">更新結果</h2>
              <div className="text-sm text-gray-200 space-y-2">
                <div>今回追加した失効エントリ数: <strong>{addedCount}</strong></div>
                <div>CRL内の失効エントリ総数: <strong>{totalRevokedCount}</strong></div>
              </div>
            </div>

            <div className="border border-gray-600 rounded p-4">
              <h2 className="text-xl font-bold mb-3 text-blue-400">更新後CRL (PEM)</h2>
              <div className="flex gap-2 mb-2">
                <button
                  onClick={() => handleCopy(crlPem, 'CRL')}
                  className="px-3 py-1 bg-green-600 hover:bg-green-500 rounded text-sm transition"
                  type="button"
                >
                  コピー
                </button>
                <button
                  onClick={() => handleDownload(crlPem, 'ca-updated.crl.pem')}
                  className="px-3 py-1 bg-blue-600 hover:bg-blue-500 rounded text-sm transition"
                  type="button"
                >
                  DL
                </button>
              </div>
              <textarea
                readOnly
                value={crlPem}
                className="w-full h-36 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs"
              />
            </div>

            {revokedSerialNumbers.length > 0 && (
              <div className="border border-gray-600 rounded p-4">
                <h2 className="text-xl font-bold mb-3 text-yellow-400">失効シリアル番号</h2>
                <div className="max-h-48 overflow-auto space-y-1 text-xs font-mono text-yellow-100">
                  {revokedSerialNumbers.map((serial, idx) => (
                    <div key={`${serial}-${idx}`}>{serial}</div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
