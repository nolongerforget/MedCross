import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import DataPage from "./pages/DataPage";
import DataUploadPage from "./pages/DataUploadPage";
import DataQueryPage from "./pages/DataQueryPage";
import ProfilePage from "./pages/ProfilePage";
import AuthPage from "./pages/AuthPage";
import BlockchainRecordsPage from "./pages/BlockchainRecordsPage";
import DataDetailPage from "./pages/DataDetailPage";
import AuthManagementPage from "./pages/AuthManagementPage";

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/data" element={<DataPage />} />
          <Route path="/data-upload" element={<DataUploadPage />} />
          <Route path="/data-query" element={<DataQueryPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/auth" element={<AuthPage />} />
          <Route path="/blockchain-records" element={<BlockchainRecordsPage />} />
          <Route path="/data-detail/:id" element={<DataDetailPage />} />
          <Route path="/auth-management" element={<AuthManagementPage />} />
          <Route path="/auth-management/:id" element={<AuthManagementPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
