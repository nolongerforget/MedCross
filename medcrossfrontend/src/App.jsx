import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import DataUploadPage from "./pages/DataUploadPage";
import DataQueryPage from "./pages/DataQueryPage";
import ProfilePage from "./pages/ProfilePage";
import AuthPage from "./pages/AuthPage";
import DataDetailPage from "./pages/DataDetailPage";

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/data-upload" element={<DataUploadPage />} />
          <Route path="/data-query" element={<DataQueryPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/auth" element={<AuthPage />} />
          <Route path="/data-detail/:id" element={<DataDetailPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
