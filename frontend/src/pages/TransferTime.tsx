// データ転送時間を計算するアプリ
import React, { useState, useEffect } from 'react';
import { CalculateTransferTime } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

interface EfficiencyResult {
  efficiency: number;
  hms: string;
  seconds: number;
  speedDisplay: string;
}

export default function TransferTime() {
  const [dataSize, setDataSize] = useSessionState('transferTime.dataSize', '');
  const [dataUnit, setDataUnit] = useSessionState('transferTime.dataUnit', 'GB');
  const [speed, setSpeed] = useSessionState('transferTime.speed', '');
  const [speedUnit, setSpeedUnit] = useSessionState('transferTime.speedUnit', 'Mbps');
  const [results, setResults] = useState<EfficiencyResult[]>([]);
  const [baseResult, setBaseResult] = useState<{ hms: string; seconds: number } | null>(null);

  useEffect(() => {
    const calculateTime = async () => {
      if (!dataSize || !speed) {
        setResults([]);
        setBaseResult(null);
        return;
      }

      const dataSizeNum = parseFloat(dataSize);
      const speedNum = parseFloat(speed);

      if (isNaN(dataSizeNum) || isNaN(speedNum) || dataSizeNum <= 0 || speedNum <= 0) {
        setResults([]);
        setBaseResult(null);
        return;
      }

      try {
        let actualSpeedUnit = speedUnit;

        const speedUnitToBps: Record<string, number> = {
          'B/s': 8,
          'KB/s': 8 * 1000,
          'MB/s': 8 * 1000 * 1000,
          'GB/s': 8 * 1000 * 1000 * 1000,
          'TB/s': 8 * 1000 * 1000 * 1000 * 1000,
          'KiB/s': 8 * 1024,
          'MiB/s': 8 * 1024 * 1024,
          'GiB/s': 8 * 1024 * 1024 * 1024,
          'TiB/s': 8 * 1024 * 1024 * 1024 * 1024,
          bps: 1,
          Kbps: 1000,
          Mbps: 1000 * 1000,
          Gbps: 1000 * 1000 * 1000,
          Tbps: 1000 * 1000 * 1000 * 1000,
          Kibps: 1024,
          Mibps: 1024 * 1024,
          Gibps: 1024 * 1024 * 1024,
          Tibps: 1024 * 1024 * 1024 * 1024
        };

        const formatSpeedDisplay = (speedInBps: number) => {
          const mbps = speedInBps / 1000 / 1000;
          const mBps = speedInBps / 8 / 1000 / 1000;
          const byteSpeed =
            mBps > 1
              ? `${mBps.toFixed(2)}MB/s`
              : `${(mBps * 1000).toFixed(2)}KB/s`;

          return `${mbps.toFixed(2)}Mbps (${byteSpeed})`;
        };

        const inputSpeedInBps = speedNum * (speedUnitToBps[speedUnit] || 0);

        // 100%効率での転送時間を計算
        const baseCalcResult = await CalculateTransferTime(dataSizeNum, dataUnit, speedNum, actualSpeedUnit);

        const formatTimeWithDays = (seconds: number) => {
          const totalSeconds = Math.floor(Number(seconds));
          const days = Math.floor(totalSeconds / 86400);
          const hours = Math.floor((totalSeconds % 86400) / 3600);
          const minutes = Math.floor((totalSeconds % 3600) / 60);
          const secs = totalSeconds % 60;

          if (days > 0) {
            return `${days}d ${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
          }
          return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
        };

        setBaseResult({
          hms: formatTimeWithDays(baseCalcResult.seconds),
          seconds: baseCalcResult.seconds
        });

        // 各効率レベルで転送時間を計算
        const efficiencyResults: EfficiencyResult[] = [];
        for (let efficiency = 100; efficiency >= 10; efficiency -= 10) {
          const adjustedSpeed = speedNum * (efficiency / 100);
          const result = await CalculateTransferTime(dataSizeNum, dataUnit, adjustedSpeed, actualSpeedUnit);
          const adjustedSpeedInBps = inputSpeedInBps * (efficiency / 100);

          efficiencyResults.push({
            efficiency,
            hms: formatTimeWithDays(result.seconds),
            seconds: result.seconds,
            speedDisplay: formatSpeedDisplay(adjustedSpeedInBps)
          });
        }

        setResults(efficiencyResults);
      } catch (error) {
        setResults([]);
        setBaseResult(null);
      }
    };

    calculateTime();
  }, [dataSize, dataUnit, speed, speedUnit]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-md">
        <h1 className="text-3xl font-bold mb-6 text-center">データ転送時間計算</h1>

        <div className="mb-4">
          <label className="block text-sm font-semibold mb-2">転送するデータ容量</label>
          <div className="flex gap-2">
            <input
              type="text"
              placeholder="100"
              value={dataSize}
              onChange={(e) => setDataSize(e.target.value)}
              className="flex-1 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            />
            <select
              name="dataUnit"
              id="dataUnit"
              value={dataUnit}
              onChange={(e) => setDataUnit(e.target.value)}
              className="p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 w-20"
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
          </div>
        </div>

        <div className="mb-6">
          <label className="block text-sm font-semibold mb-2">データ転送速度</label>
          <div className="flex gap-2">
            <input
              type="text"
              placeholder="100"
              value={speed}
              onChange={(e) => setSpeed(e.target.value)}
              className="flex-1 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            />
            <select
              name="speedUnit"
              id="speedUnit"
              value={speedUnit}
              onChange={(e) => setSpeedUnit(e.target.value)}
              className="p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 w-24"
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
          </div>
        </div>

        {baseResult && results.length > 0 && (
          <div className="mt-6">
            <div className="mb-4 p-3 bg-blue-600 rounded">
              <p className="text-sm text-gray-200">100%効率での転送時間（理論値）</p>
              <p className="text-2xl font-bold">{baseResult.hms}</p>
            </div>

            <div className="mb-3">
              <p className="text-sm font-semibold mb-3">伝送効率別転送時間</p>
              <div className="rounded border border-gray-600 overflow-hidden">
                <div className="grid grid-cols-3 bg-gray-700 px-3 py-2 text-xs font-semibold text-gray-300">
                  <span>伝送効率</span>
                  <span>転送時間</span>
                  <span>転送速度</span>
                </div>
                {results.map((item) => (
                  <div
                    key={item.efficiency}
                    className="grid grid-cols-3 px-3 py-2 border-t border-gray-700 bg-gray-800 text-sm hover:bg-gray-700 transition"
                  >
                    <p className="font-bold text-green-400">{item.efficiency}%</p>
                    <p>{item.hms}</p>
                    <p className="text-gray-200">{item.speedDisplay}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
