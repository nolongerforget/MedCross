import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Button } from "../components/ui/button";
import { Card, CardContent } from "../components/ui/card";
import { FaFileAlt, FaHistory, FaKey, FaDownload, FaChartLine } from "react-icons/fa";

/**
 * 数据详情页面
 * 显示单个医疗数据的详细信息，包括数据内容预览、元数据信息、授权历史和区块链溯源信息
 */
export default function DataDetailPage() {
  const { id } = useParams(); // 从URL获取数据ID
  const [activeTab, setActiveTab] = useState('info'); // 默认显示信息标签页
  const [dataDetail, setDataDetail] = useState(null);
  const [loading, setLoading] = useState(true);

  // 模拟从API获取数据详情
  useEffect(() => {
    // 在实际应用中，这里应该调用API获取数据
    setTimeout(() => {
      setDataDetail({
        id: id,
        fileName: '心电图数据.ecg',
        dataType: '影像数据',
        uploadTime: '2024-06-10 14:30',
        fileSize: '15.2 MB',
        owner: '张医生',
        hospital: '协和医院',
        department: '心内科',
        description: '患者王某某的心电图检查数据，采集于2024年6月10日。',
        status: '已授权',
        previewUrl: 'https://example.com/preview/123',
        // 授权历史
        authHistory: [
          { id: 1, recipient: '李医生', hospital: '人民医院', time: '2024-06-12 09:30', status: '已授权' },
          { id: 2, recipient: '医学研究中心', hospital: '医科大学', time: '2024-06-15 14:20', status: '已授权' }
        ],
        // 区块链记录
        blockchainRecords: [
          { txHash: '0x8a7d...3f21', operation: '数据上传', timestamp: '2024-06-10 14:35:22', blockNumber: 12345678 },
          { txHash: '0x9b8e...4g32', operation: '授权操作', timestamp: '2024-06-12 09:32:15', blockNumber: 12345680 },
          { txHash: '0x7c9f...5h43', operation: '授权操作', timestamp: '2024-06-15 14:22:08', blockNumber: 12345690 }
        ]
      });
      setLoading(false);
    }, 1000); // 模拟加载延迟
  }, [id]);

  if (loading) {
    return (
      <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 text-white flex items-center justify-center">
        <p className="text-xl">加载中...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 text-white">
      {/* 导航栏 */}
      <nav className="flex justify-between items-center p-6 bg-opacity-80 bg-gray-800 shadow-lg w-full">
        <h1 className="text-2xl font-bold">MedCross</h1>
        <div className="space-x-6">
          <Link to="/" className="hover:text-blue-400">首页</Link>
          <Link to="/data-upload" className="hover:text-blue-400">数据上传</Link>
          <Link to="/data-query" className="hover:text-blue-400">数据查询</Link>
          <Link to="/profile" className="hover:text-blue-400">个人中心</Link>
          <Link to="/blockchain-records" className="hover:text-blue-400">区块链记录</Link>
        </div>
      </nav>

      <div className="container mx-auto py-8 px-4">
        {/* 返回按钮 */}
        <div className="mb-6">
          <Link to="/profile" className="text-blue-400 hover:underline flex items-center">
            <span>← 返回个人中心</span>
          </Link>
        </div>

        {/* 数据标题 */}
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-3xl font-bold text-blue-400">{dataDetail.fileName}</h2>
          <div className="flex space-x-3">
            <Button className="bg-green-600 hover:bg-green-700 flex items-center">
              <FaDownload className="mr-2" /> 下载数据
            </Button>
            <Button className="bg-blue-500 hover:bg-blue-600 flex items-center">
              <FaKey className="mr-2" /> 管理授权
            </Button>
          </div>
        </div>

        {/* 标签页导航 */}
        <div className="flex border-b border-gray-700 mb-6">
          <button
            className={`py-2 px-4 ${activeTab === 'info' ? 'text-blue-400 border-b-2 border-blue-400' : 'text-gray-400 hover:text-white'}`}
            onClick={() => setActiveTab('info')}
          >
            <FaFileAlt className="inline mr-2" /> 数据信息
          </button>
          <button
            className={`py-2 px-4 ${activeTab === 'auth' ? 'text-blue-400 border-b-2 border-blue-400' : 'text-gray-400 hover:text-white'}`}
            onClick={() => setActiveTab('auth')}
          >
            <FaKey className="inline mr-2" /> 授权历史
          </button>
          <button
            className={`py-2 px-4 ${activeTab === 'blockchain' ? 'text-blue-400 border-b-2 border-blue-400' : 'text-gray-400 hover:text-white'}`}
            onClick={() => setActiveTab('blockchain')}
          >
            <FaHistory className="inline mr-2" /> 区块链记录
          </button>
        </div>

        {/* 标签页内容 */}
        <div className="bg-gray-800 rounded-xl p-6">
          {/* 数据信息标签页 */}
          {activeTab === 'info' && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              {/* 左侧数据预览 */}
              <div className="md:col-span-2">
                <h3 className="text-xl font-bold mb-4">数据预览</h3>
                <div className="bg-gray-700 rounded-lg p-4 h-64 flex items-center justify-center">
                  <div className="text-center">
                    <FaChartLine className="text-blue-400 text-5xl mx-auto mb-4" />
                    <p>心电图数据预览</p>
                    <p className="text-sm text-gray-400 mt-2">（实际应用中这里应显示数据可视化或预览）</p>
                  </div>
                </div>
                <div className="mt-6">
                  <h4 className="font-bold mb-2">数据描述</h4>
                  <p className="text-gray-300">{dataDetail.description}</p>
                </div>
              </div>
              
              {/* 右侧元数据信息 */}
              <div>
                <h3 className="text-xl font-bold mb-4">元数据信息</h3>
                <Card className="bg-gray-700">
                  <div className="space-y-3 p-4">
                    <div className="flex justify-between">
                      <span className="text-gray-400">数据ID：</span>
                      <span>{dataDetail.id}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">数据类型：</span>
                      <span>{dataDetail.dataType}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">文件大小：</span>
                      <span>{dataDetail.fileSize}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">上传时间：</span>
                      <span>{dataDetail.uploadTime}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">所有者：</span>
                      <span>{dataDetail.owner}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">所属医院：</span>
                      <span>{dataDetail.hospital}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">所属科室：</span>
                      <span>{dataDetail.department}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">授权状态：</span>
                      <span className="px-2 py-1 rounded-full text-xs bg-green-800 text-green-200">
                        {dataDetail.status}
                      </span>
                    </div>
                  </div>
                </Card>
              </div>
            </div>
          )}

          {/* 授权历史标签页 */}
          {activeTab === 'auth' && (
            <div>
              <h3 className="text-xl font-bold mb-4">授权历史</h3>
              {dataDetail.authHistory.length > 0 ? (
                <div className="bg-gray-700 rounded-lg overflow-hidden">
                  <table className="w-full">
                    <thead className="bg-gray-800">
                      <tr>
                        <th className="py-3 px-4 text-left">接收方</th>
                        <th className="py-3 px-4 text-left">所属机构</th>
                        <th className="py-3 px-4 text-left">授权时间</th>
                        <th className="py-3 px-4 text-left">状态</th>
                        <th className="py-3 px-4 text-left">操作</th>
                      </tr>
                    </thead>
                    <tbody>
                      {dataDetail.authHistory.map((auth) => (
                        <tr key={auth.id} className="border-t border-gray-600">
                          <td className="py-3 px-4">{auth.recipient}</td>
                          <td className="py-3 px-4">{auth.hospital}</td>
                          <td className="py-3 px-4">{auth.time}</td>
                          <td className="py-3 px-4">
                            <span className="px-2 py-1 rounded-full text-xs bg-green-800 text-green-200">
                              {auth.status}
                            </span>
                          </td>
                          <td className="py-3 px-4">
                            <button className="text-red-400 hover:underline">撤销授权</button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <p className="text-gray-400">暂无授权记录</p>
              )}
              
              <div className="mt-6">
                <Button className="bg-blue-500 hover:bg-blue-600">
                  <FaKey className="mr-2 inline" /> 新增授权
                </Button>
              </div>
            </div>
          )}

          {/* 区块链记录标签页 */}
          {activeTab === 'blockchain' && (
            <div>
              <h3 className="text-xl font-bold mb-4">区块链溯源记录</h3>
              <div className="bg-gray-700 rounded-lg overflow-hidden">
                <table className="w-full">
                  <thead className="bg-gray-800">
                    <tr>
                      <th className="py-3 px-4 text-left">交易哈希</th>
                      <th className="py-3 px-4 text-left">操作类型</th>
                      <th className="py-3 px-4 text-left">时间戳</th>
                      <th className="py-3 px-4 text-left">区块高度</th>
                    </tr>
                  </thead>
                  <tbody>
                    {dataDetail.blockchainRecords.map((record, index) => (
                      <tr key={index} className="border-t border-gray-600">
                        <td className="py-3 px-4 font-mono text-sm">{record.txHash}</td>
                        <td className="py-3 px-4">{record.operation}</td>
                        <td className="py-3 px-4">{record.timestamp}</td>
                        <td className="py-3 px-4">{record.blockNumber}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              
              <div className="mt-6 text-sm text-gray-400">
                <p>区块链记录不可篡改，确保数据操作的透明性和可追溯性。</p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}