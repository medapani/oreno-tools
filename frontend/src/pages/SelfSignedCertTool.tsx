import React, { useState } from 'react';
import { GenerateSelfSignedCertificate } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

export default function SelfSignedCertTool() {
  const [commonName, setCommonName] = useSessionState('selfSigned.commonName', 'localhost');
  const [organization, setOrganization] = useSessionState('selfSigned.organization', '');
  const [sanInput, setSanInput] = useSessionState('selfSigned.sanInput', 'localhost,127.0.0.1');
  const [validDays, setValidDays] = useSessionState('selfSigned.validDays', '365');
  const [algorithm, setAlgorithm] = useSessionState('selfSigned.algorithm', 'RS256');
  const [certificatePem, setCertificatePem] = useSessionState('selfSigned.certificatePem', '');
  const [privateKeyPem, setPrivateKeyPem] = useSessionState('selfSigned.privateKeyPem', '');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleGenerate = async () => {
    setError('');
    setCertificatePem('');
    setPrivateKeyPem('');

    const days = Number(validDays);
    if (!Number.isFinite(days) || days <= 0 || !Number.isInteger(days)) {
      setError('有効日数は1以上の整数で入力してください');
      return;
    }

    const sanEntries = sanInput
      .split(',')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);

    try {
      setIsLoading(true);
      const result = await GenerateSelfSignedCertificate(commonName, organization, sanEntries, days, algorithm);
      setCertificatePem(result.certificatePem);
      setPrivateKeyPem(result.privateKeyPem);
    } catch (e) {
      setError(`証明書の生成に失敗しました: ${String(e)}`);
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

  const handleClear = () => {
    setCommonName('localhost');
    setOrganization('');
    setSanInput('localhost,127.0.0.1');
    setValidDays('365');
    setAlgorithm('RS256');
    setCertificatePem('');
    setPrivateKeyPem('');
    setError('');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-5xl">
        <h1 className="text-3xl font-bold mb-6 text-center">自己署名証明書作成</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div>
            <label htmlFor="cn" className="block text-sm font-semibold mb-2">Common Name (CN)</label>
            <input
              id="cn"
              type="text"
              value={commonName}
              onChange={(e) => setCommonName(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              placeholder="localhost"
            />
          </div>
          <div>
            <label htmlFor="org" className="block text-sm font-semibold mb-2">Organization (任意)</label>
            <input
              id="org"
              type="text"
              value={organization}
              onChange={(e) => setOrganization(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              placeholder="My Org"
            />
          </div>
          <div>
            <label htmlFor="sans" className="block text-sm font-semibold mb-2">SAN (カンマ区切り)</label>
            <input
              id="sans"
              type="text"
              value={sanInput}
              onChange={(e) => setSanInput(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              placeholder="localhost,127.0.0.1,example.com"
            />
          </div>
          <div>
            <label htmlFor="days" className="block text-sm font-semibold mb-2">有効日数</label>
            <input
              id="days"
              type="number"
              min={1}
              value={validDays}
              onChange={(e) => setValidDays(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            />
          </div>
        </div>

        <div className="mb-6">
          <label htmlFor="algorithm" className="block text-sm font-semibold mb-2">鍵アルゴリズム</label>
          <select
            id="algorithm"
            value={algorithm}
            onChange={(e) => setAlgorithm(e.target.value)}
            className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
          >
            <optgroup label="RSA / RSA-PSS">
              <option value="RS256">RS256 (RSA 2048)</option>
              <option value="RS384">RS384 (RSA 2048)</option>
              <option value="RS512">RS512 (RSA 2048)</option>
              <option value="PS256">PS256 (RSA 2048)</option>
              <option value="PS384">PS384 (RSA 2048)</option>
              <option value="PS512">PS512 (RSA 2048)</option>
            </optgroup>
            <optgroup label="ECDSA">
              <option value="ES256">ES256 (P-256)</option>
              <option value="ES384">ES384 (P-384)</option>
              <option value="ES512">ES512 (P-521)</option>
            </optgroup>
            <optgroup label="EdDSA">
              <option value="EdDSA">EdDSA (Ed25519)</option>
            </optgroup>
          </select>
        </div>

        <div className="flex flex-col sm:flex-row gap-3 mb-6">
          <button
            onClick={handleGenerate}
            disabled={isLoading || !commonName.trim()}
            className="flex-1 px-4 py-2 rounded bg-blue-600 hover:bg-blue-500 disabled:bg-gray-500 disabled:cursor-not-allowed transition font-semibold"
            type="button"
          >
            {isLoading ? '生成中...' : '証明書を生成'}
          </button>
          <button
            onClick={handleClear}
            className="px-6 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
            type="button"
          >
            クリア
          </button>
        </div>

        {error && <div className="mb-6 p-3 bg-red-900 border border-red-600 rounded text-red-100">{error}</div>}

        <div className="mb-6">
          <label htmlFor="certificate-pem" className="block text-sm font-semibold mb-2">証明書 (PEM)</label>
          <textarea
            id="certificate-pem"
            value={certificatePem}
            readOnly
            className="w-full h-56 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none font-mono text-xs"
            placeholder="-----BEGIN CERTIFICATE-----"
          />
          {certificatePem && (
            <button
              onClick={() => handleCopy(certificatePem, '証明書')}
              className="mt-2 px-4 py-2 rounded bg-green-600 hover:bg-green-500 transition font-semibold text-sm"
              type="button"
            >
              証明書をコピー
            </button>
          )}
        </div>

        <div>
          <label htmlFor="private-key-pem" className="block text-sm font-semibold mb-2">秘密鍵 (PEM)</label>
          <textarea
            id="private-key-pem"
            value={privateKeyPem}
            readOnly
            className="w-full h-56 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none font-mono text-xs"
            placeholder="-----BEGIN PRIVATE KEY-----"
          />
          {privateKeyPem && (
            <button
              onClick={() => handleCopy(privateKeyPem, '秘密鍵')}
              className="mt-2 px-4 py-2 rounded bg-indigo-600 hover:bg-indigo-500 transition font-semibold text-sm"
              type="button"
            >
              秘密鍵をコピー
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
