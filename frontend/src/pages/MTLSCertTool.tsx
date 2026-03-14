import React, { useState } from 'react';
import { GenerateMTLSCertificatesMultiClient, SaveTextFile } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

interface ClientCert {
  commonName: string;
  certificatePem: string;
  privateKeyPem: string;
}

export default function MTLSCertTool() {
  const [caCommonName, setCaCommonName] = useSessionState('mtls.caCommonName', 'My Root CA');
  const [serverCommonName, setServerCommonName] = useSessionState('mtls.serverCommonName', 'localhost');
  const [clientCommonNamePrefix, setClientCommonNamePrefix] = useSessionState('mtls.clientCommonNamePrefix', 'client');
  const [clientCount, setClientCount] = useSessionState('mtls.clientCount', '3');
  const [organization, setOrganization] = useSessionState('mtls.organization', '');
  const [serverSanInput, setServerSanInput] = useSessionState('mtls.serverSanInput', 'localhost,127.0.0.1');
  const [clientSanInput, setClientSanInput] = useSessionState('mtls.clientSanInput', '');
  const [validDays, setValidDays] = useSessionState('mtls.validDays', '365');
  const [algorithm, setAlgorithm] = useSessionState('mtls.algorithm', 'RS256');

  const [caCertPem, setCaCertPem] = useSessionState('mtls.caCertPem', '');
  const [caKeyPem, setCaKeyPem] = useSessionState('mtls.caKeyPem', '');
  const [crlPem, setCrlPem] = useSessionState('mtls.crlPem', '');
  const [serverCertPem, setServerCertPem] = useSessionState('mtls.serverCertPem', '');
  const [serverKeyPem, setServerKeyPem] = useSessionState('mtls.serverKeyPem', '');
  const [clientCerts, setClientCerts] = useSessionState<ClientCert[]>('mtls.clientCerts', []);
  const [savedPaths, setSavedPaths] = useState<string[]>([]);

  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleGenerate = async () => {
    setError('');
    setCaCertPem('');
    setCaKeyPem('');
    setCrlPem('');
    setServerCertPem('');
    setServerKeyPem('');
    setClientCerts([]);
    setSavedPaths([]);

    const days = Number(validDays);
    if (!Number.isFinite(days) || days <= 0 || !Number.isInteger(days)) {
      setError('有効日数は1以上の整数で入力してください');
      return;
    }

    const count = Number(clientCount);
    if (!Number.isInteger(count) || count <= 0 || count > 100) {
      setError('クライアント証明書の数は1から100の間で指定してください');
      return;
    }

    const serverSanEntries = serverSanInput
      .split(',')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);

    const clientSanEntries = clientSanInput
      .split(',')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);

    try {
      setIsLoading(true);
      const result = await GenerateMTLSCertificatesMultiClient(
        caCommonName,
        serverCommonName,
        clientCommonNamePrefix,
        count,
        organization,
        serverSanEntries,
        clientSanEntries,
        days,
        algorithm
      );

      setCaCertPem(result.caCertificatePem);
      setCaKeyPem(result.caPrivateKeyPem);
      setCrlPem(result.crlPem);
      setServerCertPem(result.serverCertificatePem);
      setServerKeyPem(result.serverPrivateKeyPem);
      setClientCerts(result.clientCertificates || []);
    } catch (e) {
      setError(`mTLS証明書の生成に失敗しました: ${String(e)}`);
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

  const handleDownloadAll = async () => {
    const downloaded: string[] = [];

    // CA証明書
    downloaded.push(await handleDownload(caCertPem, 'ca-cert.pem', false));
    downloaded.push(await handleDownload(caKeyPem, 'ca-key.pem', false));
    downloaded.push(await handleDownload(crlPem, 'ca.crl.pem', false));

    // サーバー証明書
    downloaded.push(await handleDownload(serverCertPem, 'server-cert.pem', false));
    downloaded.push(await handleDownload(serverKeyPem, 'server-key.pem', false));

    // クライアント証明書
    for (const client of clientCerts) {
      downloaded.push(await handleDownload(client.certificatePem, `${client.commonName}-cert.pem`, false));
      downloaded.push(await handleDownload(client.privateKeyPem, `${client.commonName}-key.pem`, false));
    }

    const validPaths = downloaded.filter((p) => p);
    if (validPaths.length > 0) {
      setSavedPaths(validPaths);
    }

    alert(`保存完了: ~/Downloads/oreno-tools-certs に全${3 + 2 + clientCerts.length * 2}個のファイルを出力しました`);
  };

  const handleClear = () => {
    setCaCommonName('My Root CA');
    setServerCommonName('localhost');
    setClientCommonNamePrefix('client');
    setClientCount('3');
    setOrganization('');
    setServerSanInput('localhost,127.0.0.1');
    setClientSanInput('');
    setValidDays('365');
    setAlgorithm('RS256');
    setCaCertPem('');
    setCaKeyPem('');
    setCrlPem('');
    setServerCertPem('');
    setServerKeyPem('');
    setClientCerts([]);
    setSavedPaths([]);
    setError('');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-6xl">
        <h1 className="text-3xl font-bold mb-6 text-center">mTLS 証明書作成</h1>
        <p className="text-gray-400 text-center mb-6">CA、サーバー、複数のクライアント証明書を一度に生成します</p>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div>
            <label htmlFor="ca-cn" className="block text-sm font-semibold mb-2">CA Common Name</label>
            <input
              id="ca-cn"
              type="text"
              value={caCommonName}
              onChange={(e) => setCaCommonName(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              placeholder="My Root CA"
            />
          </div>
          <div>
            <label htmlFor="server-cn" className="block text-sm font-semibold mb-2">Server Common Name</label>
            <input
              id="server-cn"
              type="text"
              value={serverCommonName}
              onChange={(e) => setServerCommonName(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              placeholder="localhost"
            />
          </div>
          <div>
            <label htmlFor="client-prefix" className="block text-sm font-semibold mb-2">クライアント証明書名プレフィックス</label>
            <input
              id="client-prefix"
              type="text"
              value={clientCommonNamePrefix}
              onChange={(e) => setClientCommonNamePrefix(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              placeholder="client"
            />
          </div>
          <div>
            <label htmlFor="client-count" className="block text-sm font-semibold mb-2">クライアント証明書の数（1-100）</label>
            <input
              id="client-count"
              type="number"
              min={1}
              max={100}
              value={clientCount}
              onChange={(e) => setClientCount(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
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
              placeholder="My Organization"
            />
          </div>
          <div>
            <label htmlFor="server-san" className="block text-sm font-semibold mb-2">Server SAN (カンマ区切り)</label>
            <input
              id="server-san"
              type="text"
              value={serverSanInput}
              onChange={(e) => setServerSanInput(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              placeholder="localhost,127.0.0.1,example.com"
            />
          </div>
          <div>
            <label htmlFor="client-san" className="block text-sm font-semibold mb-2">Client SAN (任意、カンマ区切り)</label>
            <input
              id="client-san"
              type="text"
              value={clientSanInput}
              onChange={(e) => setClientSanInput(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              placeholder="client.example.com"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
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
          <div>
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
        </div>

        <div className="flex flex-col gap-3 mb-6">
          <button
            onClick={handleGenerate}
            disabled={isLoading || !caCommonName.trim() || !serverCommonName.trim() || !clientCommonNamePrefix.trim()}
            className="w-full px-4 py-3 rounded bg-blue-600 hover:bg-blue-500 disabled:bg-gray-500 disabled:cursor-not-allowed transition font-semibold text-lg"
            type="button"
          >
            {isLoading ? '生成中...' : '証明書を生成'}
          </button>
          {caCertPem && (
            <button
              onClick={handleDownloadAll}
              className="w-full px-4 py-3 rounded bg-green-600 hover:bg-green-500 transition font-semibold text-lg"
              type="button"
            >
              すべてダウンロード
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

        {caCertPem && (
          <div className="space-y-6">
            <div className="border border-gray-600 rounded p-4">
              <h2 className="text-xl font-bold mb-3 text-blue-400">CA 証明書</h2>
              <div className="space-y-3">
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <label className="text-sm font-semibold text-gray-300">証明書 (PEM)</label>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleCopy(caCertPem, 'CA証明書')}
                        className="px-3 py-1 bg-green-600 hover:bg-green-500 rounded text-sm transition"
                        type="button"
                      >
                        コピー
                      </button>
                      <button
                        onClick={() => handleDownload(caCertPem, 'ca-cert.pem')}
                        className="px-3 py-1 bg-blue-600 hover:bg-blue-500 rounded text-sm transition"
                        type="button"
                      >
                        DL
                      </button>
                    </div>
                  </div>
                  <textarea
                    readOnly
                    value={caCertPem}
                    className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs"
                  />
                </div>
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <label className="text-sm font-semibold text-gray-300">秘密鍵 (PEM)</label>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleCopy(caKeyPem, 'CA秘密鍵')}
                        className="px-3 py-1 bg-green-600 hover:bg-green-500 rounded text-sm transition"
                        type="button"
                      >
                        コピー
                      </button>
                      <button
                        onClick={() => handleDownload(caKeyPem, 'ca-key.pem')}
                        className="px-3 py-1 bg-blue-600 hover:bg-blue-500 rounded text-sm transition"
                        type="button"
                      >
                        DL
                      </button>
                    </div>
                  </div>
                  <textarea
                    readOnly
                    value={caKeyPem}
                    className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs"
                  />
                </div>
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <label className="text-sm font-semibold text-gray-300">CRL (PEM)</label>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleCopy(crlPem, 'CRL')}
                        className="px-3 py-1 bg-green-600 hover:bg-green-500 rounded text-sm transition"
                        type="button"
                      >
                        コピー
                      </button>
                      <button
                        onClick={() => handleDownload(crlPem, 'ca.crl.pem')}
                        className="px-3 py-1 bg-blue-600 hover:bg-blue-500 rounded text-sm transition"
                        type="button"
                      >
                        DL
                      </button>
                    </div>
                  </div>
                  <textarea
                    readOnly
                    value={crlPem}
                    className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs"
                  />
                </div>
              </div>
            </div>

            <div className="border border-gray-600 rounded p-4">
              <h2 className="text-xl font-bold mb-3 text-green-400">サーバー証明書</h2>
              <div className="space-y-3">
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <label className="text-sm font-semibold text-gray-300">証明書 (PEM)</label>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleCopy(serverCertPem, 'サーバー証明書')}
                        className="px-3 py-1 bg-green-600 hover:bg-green-500 rounded text-sm transition"
                        type="button"
                      >
                        コピー
                      </button>
                      <button
                        onClick={() => handleDownload(serverCertPem, 'server-cert.pem')}
                        className="px-3 py-1 bg-blue-600 hover:bg-blue-500 rounded text-sm transition"
                        type="button"
                      >
                        DL
                      </button>
                    </div>
                  </div>
                  <textarea
                    readOnly
                    value={serverCertPem}
                    className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs"
                  />
                </div>
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <label className="text-sm font-semibold text-gray-300">秘密鍵 (PEM)</label>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleCopy(serverKeyPem, 'サーバー秘密鍵')}
                        className="px-3 py-1 bg-green-600 hover:bg-green-500 rounded text-sm transition"
                        type="button"
                      >
                        コピー
                      </button>
                      <button
                        onClick={() => handleDownload(serverKeyPem, 'server-key.pem')}
                        className="px-3 py-1 bg-blue-600 hover:bg-blue-500 rounded text-sm transition"
                        type="button"
                      >
                        DL
                      </button>
                    </div>
                  </div>
                  <textarea
                    readOnly
                    value={serverKeyPem}
                    className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs"
                  />
                </div>
              </div>
            </div>

            <div className="border border-gray-600 rounded p-4">
              <h2 className="text-xl font-bold mb-3 text-yellow-400">クライアント証明書 ({clientCerts.length}個)</h2>
              {clientCerts.map((client, index) => (
                <div key={index} className="mb-6 pb-6 border-b border-gray-700 last:border-b-0 last:mb-0 last:pb-0">
                  <h3 className="text-lg font-semibold mb-3 text-yellow-300 border-l-4 border-yellow-500 pl-3">{client.commonName}</h3>
                  <div className="space-y-3">
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <label className="text-sm font-semibold text-gray-300">証明書 (PEM)</label>
                        <div className="flex gap-2">
                          <button
                            onClick={() => handleCopy(client.certificatePem, `${client.commonName}証明書`)}
                            className="px-3 py-1 bg-green-600 hover:bg-green-500 rounded text-sm transition"
                            type="button"
                          >
                            コピー
                          </button>
                          <button
                            onClick={() => handleDownload(client.certificatePem, `${client.commonName}-cert.pem`)}
                            className="px-3 py-1 bg-blue-600 hover:bg-blue-500 rounded text-sm transition"
                            type="button"
                          >
                            DL
                          </button>
                        </div>
                      </div>
                      <textarea
                        readOnly
                        value={client.certificatePem}
                        className="w-full h-24 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs"
                      />
                    </div>
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <label className="text-sm font-semibold text-gray-300">秘密鍵 (PEM)</label>
                        <div className="flex gap-2">
                          <button
                            onClick={() => handleCopy(client.privateKeyPem, `${client.commonName}秘密鍵`)}
                            className="px-3 py-1 bg-green-600 hover:bg-green-500 rounded text-sm transition"
                            type="button"
                          >
                            コピー
                          </button>
                          <button
                            onClick={() => handleDownload(client.privateKeyPem, `${client.commonName}-key.pem`)}
                            className="px-3 py-1 bg-blue-600 hover:bg-blue-500 rounded text-sm transition"
                            type="button"
                          >
                            DL
                          </button>
                        </div>
                      </div>
                      <textarea
                        readOnly
                        value={client.privateKeyPem}
                        className="w-full h-24 p-3 rounded bg-gray-700 border border-gray-600 font-mono text-xs"
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
