import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Card } from "../components/ui/card";
import { FaSearch, FaExchangeAlt, FaFileAlt, FaKey, FaInfoCircle } from "react-icons/fa";

/**
 * 区块链交易记录页面
 * 显示与用户相关的所有区块链交易记录，包括数据上传、授权操作等
 * 提供完整的数据操作溯源功能
 */
export default function BlockchainRecordsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState('all'); // 筛选类型：all, upload, auth
  const [records, setRecords] = useState([]);
  const [loading, setLoading] = useState(true);

  // 模拟从API获取区块链记录
  useEffect(() => {
    setTimeout(() => {
      const mockRecords = [
        {
          txHash: '0x8a7d...3f21',
          operation: '数据上传',
          dataName: '心电图数据.ecg',
          dataType: '影像数据',
          timestamp: '2024-06-10 14:35:22',
          blockNumber: 12345678,
          status: '已确认',
          dataId: '1'
        },
        {
          txHash: '0x9b8e...4g32',
          operation: '授权操作',
          dataName: '心电图数据.ecg',
          dataType: '影像数据',
          recipient: '李医生',
          recipientOrg: '人民医院',
          timestamp: '2024-06-12 09:32:15',
          blockNumber: 12345680,
          status: '已确认',
          dataId: '1'
        },
        {
          txHash: '0x7c9f...5h43',
          operation: '授权操作',
          dataName: '基因测序结果.xml',
          dataType: '基因组数据',
          recipient: '医学研究中心',
          recipientOrg: '医科大学',
          timestamp: '2024-06-15 14:22:08',
          blockNumber: 12345690,
          status: '已确认',
          dataId: '3'
        },
        {
          txHash: '0x6d8e...7j54',
          operation: '数据上传',
          dataName: '患者病历记录.pdf',
          dataType: '电子病历',
          timestamp: '2024-06-08 09:20:33',
          blockNumber: 12345670,
          status: '已确认',
          dataId: '2'
        },
        {
          txHash: '0x5e7f...8k65',
          operation: '数据访问',
          dataName: '心电图数据.ecg',
          dataType: '影像数据',
          accessor: '李医生',
          accessorOrg: '人民医院',
          timestamp: '2024-06-13 10:45:19',
          blockNumber: 12345685,
          status: '已确认',
          dataId: '1'
        }
      ];
      setRecords(mockRecords);
      setLoading(false);
    }, 1000); // 模拟加载延迟
  }, []);

  // 处理搜索和筛选
  const getFilteredRecords = () => {
    let filtered = records;
    
    // 按类型筛选
    if (filterType !== 'all') {
      filtered = filtered.filter(record => {
        if (filterType === 'upload') return record.operation === '数据上传';
        if (filterType === 'auth') return record.operation === '授权操作';
        if (filterType === 'access') return record.operation === '数据访问';
        return true;
      });
    }
    
    // 按搜索词筛选
    if (searchTerm) {
      filtered = filtered.filter(record => 
        record.txHash.includes(searchTerm) ||
        record.dataName.includes(searchTerm) ||
        (record.recipient && record.recipient.includes(searchTerm)) ||
        (record.accessor && record.accessor.includes(searchTerm))
      );
    }
    
    return filtered;
  };

  // 获取操作图标
  const getOperationIcon = (operation) => {
    switch (operation) {
      case '数据上传': return <FaFileAlt className="text-blue-400" />;
      case '授权操作': return <FaKey className="text-green-400" />;
      case '数据访问': return <FaSearch className="text-yellow-400" />;
      default: return <FaExchangeAlt className="text-gray-400" />;
    }
  };

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
          <Link to="/blockchain-records" className="text-blue-400">区块链记录</Link>
        </div>
      </nav>

      <div className="container mx-auto py-8 px-4">
        {/* 返回按钮 */}
        <div className="mb-6">
          <Link to="/profile" className="text-blue-400 hover:underline flex items-center">
            <span>← 返回个人中心</span>
          </Link>
        </div>

        <h2 className="text-3xl font-bold text-blue-400 mb-6">区块链交易记录</h2>

        {/* 搜索和筛选 */}
        <div className="mb-6 flex flex-wrap gap-4">
          <div className="flex-grow">
            <Input
              type="text"
              placeholder="搜索交易哈希或数据名称"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full bg-gray-700 border-gray-600"
            />
          </div>
          
          <div className="flex space-x-2">
            <Button 
              onClick={() => setFilterType('all')} 
              className={`${filterType === 'all' ? 'bg-blue-500' : 'bg-gray-700'} hover:bg-blue-600`}
            >
              全部
            </Button>
            <Button 
              onClick={() => setFilterType('upload')} 
              className={`${filterType === 'upload' ? 'bg-blue-500' : 'bg-gray-700'} hover:bg-blue-600`}
            >
              <FaFileAlt className="mr-2" /> 上传
            </Button>
            <Button 
              onClick={() => setFilterType('auth')} 
              className={`${filterType === 'auth' ? 'bg-blue-500' : 'bg-gray-700'} hover:bg-blue-600`}
            >
              <FaKey className="mr-2" /> 授权
            </Button>
            <Button 
              onClick={() => setFilterType('access')} 
              className={`${filterType === 'access' ? 'bg-blue-500' : 'bg-gray-700'} hover:bg-blue-600`}
            >
              <FaSearch className="mr-2" /> 访问
            </Button>
          </div>
        </div>

        {/* 区块链记录列表 */}
        <div className="space-y-4">
          {getFilteredRecords().length > 0 ? (
            getFilteredRecords().map((record, index) => (
              <Card key={index} className="bg-gray-800 p-4 hover:bg-gray-750 transition-colors">
                <div className="flex items-start">
                  <div className="p-3 bg-gray-700 rounded-full mr-4">
                    {getOperationIcon(record.operation)}
                  </div>
                  
                  <div className="flex-grow">
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="font-bold text-lg">{record.operation}</h3>
                        <p className="text-gray-400 text-sm">
                          交易哈希: <span className="font-mono">{record.txHash}</span>
                        </p>
                      </div>
                      <span className="px-2 py-1 rounded-full text-xs bg-green-800 text-green-200">
                        {record.status}
                      </span>
                    </div>
                    
                    <div className="mt-2 grid grid-cols-1 md:grid-cols-2 gap-2">
                      <div>
                        <p>
                          <span className="text-gray-400">数据名称: </span>
                          <Link to={`/data-detail/${record.dataId}`} className="text-blue-400 hover:underline">
                            {record.dataName}
                          </Link>
                        </p>
                        <p><span className="text-gray-400">数据类型: </span>{record.dataType}</p>
                        {record.recipient && (
                          <p><span className="text-gray-400">接收方: </span>{record.recipient} ({record.recipientOrg})</p>
                        )}
                        {record.accessor && (
                          <p><span className="text-gray-400">访问者: </span>{record.accessor} ({record.accessorOrg})</p>
                        )}
                      </div>
                      <div>
                        <p><span className="text-gray-400">时间戳: </span>{record.timestamp}</p>
                        <p><span className="text-gray-400">区块高度: </span>{record.blockNumber}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </Card>
            ))
          ) : (
            <div className="text-center py-10">
              <FaInfoCircle className="text-5xl text-gray-600 mx-auto mb-4" />
              <p className="text-xl text-gray-400">没有找到匹配的交易记录</p>
            </div>
          )}
        </div>

        {/* 区块链信息说明 */}
        <div className="mt-8 bg-gray-800 p-4 rounded-lg">
          <h3 className="font-bold mb-2 flex items-center">
            <FaInfoCircle className="mr-2 text-blue-400" /> 关于区块链记录
          </h3>
          <p className="text-sm text-gray-400">
            区块链记录提供了所有数据操作的不可篡改证明。每一条记录都被永久存储在区块链上，
            确保数据的完整性和可追溯性。通过这些记录，您可以清晰地了解数据的流转过程，
            包括谁上传了数据、谁被授权访问以及何时进行了这些操作。
          </p>
        </div>
      </div>
    </div>
  );
}