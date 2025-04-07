import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Card, CardContent } from "../components/ui/card";
import { FaCloudUploadAlt, FaFileAlt, FaTimesCircle, FaNetworkWired } from "react-icons/fa";
import { FaEthereum } from "react-icons/fa6";

/**
 * 数据上传页面
 * 用于医疗数据的上传，包括文件选择、元数据填写和上传进度显示
 * 支持多种医疗数据类型的上传，并提供实时上传进度反馈
 */
export default function DataUploadPage() {
  // 文件上传状态管理
  const [file, setFile] = useState(null);
  const [dataType, setDataType] = useState('');
  const [description, setDescription] = useState('');
  const [tags, setTags] = useState('');
  const [targetChain, setTargetChain] = useState('ethereum'); // 默认选择以太坊链
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadSuccess, setUploadSuccess] = useState(false);
  
  // 处理文件选择
  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile) {
      setFile(selectedFile);
      // 根据文件类型自动选择数据类型
      const fileExt = selectedFile.name.split('.').pop().toLowerCase();
      if (['jpg', 'png', 'dcm', 'dicom'].includes(fileExt)) {
        setDataType('影像数据');
      } else if (['pdf', 'doc', 'docx'].includes(fileExt)) {
        setDataType('电子病历');
      } else if (['xml', 'json', 'csv'].includes(fileExt)) {
        setDataType('基因组数据');
      }
    }
  };

  // 清除已选文件
  const clearFile = () => {
    setFile(null);
    setUploadProgress(0);
    setUploadSuccess(false);
  };

  // 处理上传操作
  const handleUpload = () => {
    if (!file || !dataType) {
      alert("请选择文件并选择数据类型！");
      return;
    }

    if (!description) {
      alert("请填写数据描述！");
      return;
    }
    
    if (!targetChain) {
      alert("请选择目标区块链！");
      return;
    }

    // 开始上传
    setIsUploading(true);
    setUploadProgress(0);

    // 模拟上传进度
    const interval = setInterval(() => {
      setUploadProgress(prev => {
        if (prev >= 100) {
          clearInterval(interval);
          setIsUploading(false);
          setUploadSuccess(true);
          return 100;
        }
        return prev + 10;
      });
    }, 500);

    // 模拟上传过程（将文件数据和类型发送给后端）
    console.log("上传文件:", file);
    console.log("数据类型:", dataType);
    console.log("数据描述:", description);
    console.log("关键词:", tags);
    console.log("目标区块链:", targetChain);
  };

  return (
    <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 text-white">
      {/* 导航栏 */}
      <nav className="flex justify-between items-center p-6 bg-opacity-80 bg-gray-800 shadow-lg w-full">
        <h1 className="text-2xl font-bold">MedCross</h1>
        <div className="space-x-6">
          <Link to="/" className="hover:text-blue-400">首页</Link>
          <Link to="/data-upload" className="text-blue-400">数据上传</Link>
          <Link to="/data-query" className="hover:text-blue-400">数据查询</Link>
          <Link to="/profile" className="hover:text-blue-400">个人中心</Link>
        </div>
      </nav>

      <div className="container mx-auto py-8 px-4">
        <h2 className="text-3xl font-bold text-blue-400 mb-6">医疗数据上传</h2>
        
        <Card className="bg-gray-800 p-6 rounded-xl shadow-lg">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            {/* 左侧文件上传区域 */}
            <div>
              <h3 className="text-xl font-bold mb-4">选择文件</h3>
              
              {!file ? (
                <div className="border-2 border-dashed border-gray-600 rounded-lg p-8 text-center cursor-pointer hover:border-blue-400 transition-colors"
                     onClick={() => document.getElementById('fileInput').click()}>
                  <FaCloudUploadAlt className="text-5xl text-blue-400 mx-auto mb-4" />
                  <p className="mb-2">点击或拖拽文件到此处上传</p>
                  <p className="text-sm text-gray-400">支持的文件类型：影像数据、电子病历、基因组数据等</p>
                  <input
                    id="fileInput"
                    type="file"
                    onChange={handleFileChange}
                    className="hidden"
                  />
                </div>
              ) : (
                <div className="bg-gray-700 rounded-lg p-4">
                  <div className="flex justify-between items-center mb-2">
                    <div className="flex items-center">
                      <FaFileAlt className="text-blue-400 mr-2" />
                      <span className="truncate">{file.name}</span>
                    </div>
                    <button onClick={clearFile} className="text-gray-400 hover:text-red-400">
                      <FaTimesCircle />
                    </button>
                  </div>
                  <div className="text-sm text-gray-400 mb-2">{(file.size / 1024 / 1024).toFixed(2)} MB</div>
                  
                  {isUploading && (
                    <div className="w-full bg-gray-600 rounded-full h-2 mb-2">
                      <div 
                        className="bg-blue-500 h-2 rounded-full transition-all duration-300" 
                        style={{ width: `${uploadProgress}%` }}
                      ></div>
                    </div>
                  )}
                  
                  {uploadSuccess && (
                    <div className="text-green-400 text-sm mb-2">上传成功！</div>
                  )}
                </div>
              )}
            </div>
            
            {/* 选择目标区块链 */}
            <div className="mt-6">
              <h3 className="text-xl font-bold mb-4">选择目标区块链</h3>
              <div className="grid grid-cols-2 gap-4">
                <div 
                  className={`p-4 rounded-lg cursor-pointer flex items-center ${targetChain === 'ethereum' ? 'bg-blue-800 border-2 border-blue-400' : 'bg-gray-700 hover:bg-gray-600'}`}
                  onClick={() => setTargetChain('ethereum')}
                >
                  <FaEthereum className="text-3xl text-blue-400 mr-3" />
                  <div>
                    <p className="font-medium">以太坊</p>
                    <p className="text-sm text-gray-400">适合元数据和访问控制</p>
                  </div>
                </div>
                <div 
                  className={`p-4 rounded-lg cursor-pointer flex items-center ${targetChain === 'fabric' ? 'bg-blue-800 border-2 border-blue-400' : 'bg-gray-700 hover:bg-gray-600'}`}
                  onClick={() => setTargetChain('fabric')}
                >
                  <FaNetworkWired className="text-3xl text-blue-400 mr-3" />
                  <div>
                    <p className="font-medium">Hyperledger Fabric</p>
                    <p className="text-sm text-gray-400">适合详细医疗数据记录</p>
                  </div>
                </div>
              </div>
            </div>
            
            {/* 右侧元数据填写区域 */}
            <div>
              <h3 className="text-xl font-bold mb-4">数据信息</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-1">数据类型</label>
                  <select
                    className="w-full p-2 rounded-md bg-gray-700 border-gray-600"
                    value={dataType}
                    onChange={(e) => setDataType(e.target.value)}
                  >
                    <option value="">选择数据类型</option>
                    <option value="电子病历">电子病历</option>
                    <option value="基因组数据">基因组数据</option>
                    <option value="影像数据">影像数据</option>
                    <option value="检验报告">检验报告</option>
                    <option value="处方数据">处方数据</option>
                  </select>
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">数据描述</label>
                  <textarea
                    className="w-full p-2 rounded-md bg-gray-700 border-gray-600 min-h-[100px]"
                    placeholder="请简要描述该数据的内容、来源和用途"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                  ></textarea>
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">关键词（用逗号分隔）</label>
                  <Input
                    type="text"
                    placeholder="例如：心脏病,急诊,2024年"
                    value={tags}
                    onChange={(e) => setTags(e.target.value)}
                    className="w-full bg-gray-700 border-gray-600"
                  />
                  <p className="text-sm text-gray-400 mt-1">关键词将用于跨链数据查询</p>
                </div>
              </div>
            </div>
          </div>
          
          <div className="mt-8 flex justify-center">
            <Button 
              onClick={handleUpload} 
              disabled={!file || isUploading}
              className="bg-blue-500 hover:bg-blue-600 px-8 py-2 text-lg disabled:bg-gray-600 disabled:cursor-not-allowed"
            >
              {isUploading ? '上传中...' : `上传到${targetChain === 'ethereum' ? '以太坊' : 'Fabric'}链`}
            </Button>
          </div>
        </Card>
        
        {/* 上传指南 */}
        <Card className="bg-gray-800 p-6 rounded-xl shadow-lg mt-8">
          <h3 className="text-xl font-bold mb-4">数据上传指南</h3>
          <div className="text-gray-300 space-y-2">
            <p>1. 请确保您有权上传和分享该医疗数据。</p>
            <p>2. 上传前请移除所有可能的患者身份标识信息，除非您已获得适当授权。</p>
            <p>3. 所有上传的数据将通过区块链技术进行记录，确保数据溯源和不可篡改性。</p>
            <p>4. 您可以在上传后通过授权管理页面控制谁可以访问您的数据。</p>
            <p>5. 不同区块链适用于不同场景：以太坊适合元数据和访问控制，Fabric适合详细医疗数据。</p>
          </div>
        </Card>
      </div>
    </div>
  );
}