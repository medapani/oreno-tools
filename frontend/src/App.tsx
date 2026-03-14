import { HashRouter, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import NetworkUnit from './pages/NetworkUnit';
import DataTransferRate from './pages/DataTransferRate';
import TransferTime from './pages/TransferTime';
import Base64Tool from './pages/Base64Tool';
import JWTTool from './pages/JWTTool';
import JSONFormatter from './pages/JSONFormatter';
import URLTool from './pages/URLTool';
import KeyPairTool from './pages/KeyPairTool';
import SelfSignedCertTool from './pages/SelfSignedCertTool';
import MTLSCertTool from './pages/MTLSCertTool';
import CRLTool from './pages/CRLTool';
import CIDRCalculator from './pages/CIDRCalculator';
import BaseConverter from './pages/BaseConverter';
import Layout from './components/Layout';

function App() {
    return (
        <HashRouter>
            <Routes>
                <Route element={<Layout />}>
                    <Route path="/" element={<Home />} />
                    <Route path="/network-unit" element={<NetworkUnit />} />
                    <Route path="/data-transfer-rate" element={<DataTransferRate />} />
                    <Route path="/transfer-time" element={<TransferTime />} />
                    <Route path="/base64" element={<Base64Tool />} />
                    <Route path="/jwt" element={<JWTTool />} />
                    <Route path="/json-formatter" element={<JSONFormatter />} />
                    <Route path="/url-tool" element={<URLTool />} />
                    <Route path="/key-pair" element={<KeyPairTool />} />
                    <Route path="/self-signed-cert" element={<SelfSignedCertTool />} />
                    <Route path="/mtls-cert" element={<MTLSCertTool />} />
                    <Route path="/crl-tool" element={<CRLTool />} />
                    <Route path="/cidr-calculator" element={<CIDRCalculator />} />
                    <Route path="/base-converter" element={<BaseConverter />} />
                </Route>
            </Routes>
        </HashRouter>
    );
}

export default App;