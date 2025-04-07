import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Card, CardContent } from "../components/ui/card";
import { FaSearch, FaFilter, FaFileAlt, FaDna, FaImage, FaFileMedical, FaEye, FaNetworkWired } from "react-icons/fa";
import { FaEthereum } from "react-icons/fa6";

/**
 * 数据查询页面
 * 用于医疗数据的检索和查询，提供多种筛选条件和结果展示
 * 支持按数据类型、上传时间、关键词等多维度查询
 */
export default function DataQueryPage() {
  // 查询状态管理
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedType, setSelectedType] = useState('all');
  const [dateRange, setDateRange] = useState({ start: '', end: '' });
  const [sortBy, setSortBy] = useState('newest');
  const [showFilters, setShowFilters] = useState(false);
  const [searchResults, setSearchResults] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [chainSource, setChainSource] = useState('all'); // 区块链源筛选
  
  // 模拟从API获取数据
  useEffect(() => {
    // 初始加载一些示例数据
    const mockData = [
      {
        id: '1',
        fileName: '心电图数据.ecg',
        dataType: '影像数据',
        uploadTime: '2024-06-10 14:30',
        fileSize: '15.2 MB',
        owner: '张医生',
        hospital: '协和医院',
        description: '患者王某某的心电图检查数据，采集于2024年6月10日。',
        tags: ['心脏病', '急诊', '2024年'],
        chain: 'ethereum'
      },
      {
        id: '2',
        fileName: '患者病历记录.pdf',
        dataType: '电子病历',
        uploadTime: '2024-06-08 09:15',
        fileSize: '2.8 MB',
        owner: '李医生',
        hospital: '人民医院',
        description: '患者张某某的完整病历记录，包含既往病史和当前治疗方案。',
        tags: ['内科', '慢性病', '随访'],
        chain: 'fabric'
      },
      {
        id: '3',
        fileName: '基因测序结果.xml',
        dataType: '基因组数据',
        uploadTime: '2024-06-01 16:45',
        fileSize: '128.5 MB',
        owner: '王研究员',
        hospital: '医学研究中心',
        description: '肿瘤患者基因测序数据，用于精准医疗研究。',
        tags: ['肿瘤', '基因测序', '精准医疗'],
        chain: 'ethereum'
      },
      {
        id: '4',
        fileName: 'CT扫描图像.dicom',
        dataType: '影像数据',
        uploadTime: '2024-06-05 11:20',
        fileSize: '45.7 MB',
        owner: '赵医生',
        hospital: '第三医院',
        description: '肺部CT扫描图像，用于肺炎诊断。',
        tags: ['肺炎', 'CT', '影像学'],
        chain: 'fabric'
      },
      {
        id: '5',
        fileName: '处方数据集.csv',
        dataType: '处方数据',
        uploadTime: '2024-06-12 10:05',
        fileSize: '8.3 MB',
        owner: '刘药师',
        hospital: '中心医院',
        description: '2024年5月门诊处方数据汇总，用于药物使用分析。',
        tags: ['处方', '药物', '统计分析'],
        chain: 'ethereum'
      },
    ];
    
    setSearchResults(mockData);
  }, []);

  // 处理查询
  const handleSearch = () => {
    setIsLoading(true);
    
    // 模拟API请求延迟
    setTimeout(() => {
      // 在实际应用中，这里应该调用后端API进行查询
      // 请求参数应包含：关键词、数据类型、日期范围、排序方式、区块链源
      console.log("执行跨链查询:", {
        keyword: searchQuery,
        dataType: selectedType,
        dateRange: dateRange,
        sortBy: sortBy,
        chainSource: chainSource
      });
      
      // 模拟从API获取数据
      let filteredResults = [...searchResults];
      
      // 按关键词筛选
      if (searchQuery) {
        filteredResults = filteredResults.filter(item =>
          item.fileName.toLowerCase().includes(searchQuery.toLowerCase()) ||
          item.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
          item.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()))
        );
      }
      
      // 按数据类型筛选
      if (selectedType !== 'all') {
        filteredResults = filteredResults.filter(item => item.dataType === selectedType);
      }
      
      // 按区块链源筛选
      if (chainSource !== 'all') {
        filteredResults = filteredResults.filter(item => item.chain === chainSource);
      }
      
      // 按日期范围筛选
      if (dateRange.start) {
        const startDate = new Date(dateRange.start);
        filteredResults = filteredResults.filter(item => new Date(item.uploadTime) >= startDate);
      }
      
      if (dateRange.end) {
        const endDate = new Date(dateRange.end);
        endDate.setHours(23, 59, 59); // 设置为当天结束时间
        filteredResults = filteredResults.filter(item => new Date(item.uploadTime) <= endDate);
      }
      
      // 排序
      if (sortBy === 'newest') {
        filteredResults.sort((a, b) => new Date(b.uploadTime) - new Date(a.uploadTime));
      } else if (sortBy === 'oldest') {
        filteredResults.sort((a, b) => new Date(a.uploadTime) - new Date(b.uploadTime));
      } else if (sortBy === 'name') {
        filteredResults.sort((a, b) => a.fileName.localeCompare(b.fileName));
      }
      
      setSearchResults(filteredResults);
      setIsLoading(false);
    }, 800);
  };

  // 获取数据类型图标
  const getDataTypeIcon = (type) => {
    switch (type) {
      case '电子病历': return <FaFileMedical className="text-blue-400" />;
      case '基因组数据': return <FaDna className="text-green-400" />;
      case '影像数据': return <FaImage className="text-purple-400" />;
      case '处方数据': return <FaFileAlt className="text-yellow-400" />;
      default: return <FaFileAlt className="text-gray-400" />;
    }
  };

  // 重置筛选条件
  const resetFilters = () => {
    setSearchQuery('');
    setSelectedType('all');
    setDateRange({ start: '', end: '' });
    setSortBy('newest');
    setChainSource('all'); // 重置区块链源筛选
  };

  return (
    <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 text-white">
      {/* 导航栏 */}
      <nav className="flex justify-between items-center p-6 bg-opacity-80 bg-gray-800 shadow-lg w-full">
        <h1 className="text-2xl font-bold">MedCross</h1>
        <div className="space-x-6">
          <Link to="/" className="hover:text-blue-400">首页</Link>
          <Link to="/data-upload" className="hover:text-blue-400">数据上传</Link>
          <Link to="/data-query" className="text-blue-400">数据查询</Link>
          <Link to="/profile" className="hover:text-blue-400">个人中心</Link>
        </div>
      </nav>

      <div className="container mx-auto py-8 px-4">
        <h2 className="text-3xl font-bold text-blue-400 mb-6">医疗数据查询</h2>
        
        {/* 搜索栏 */}
        <Card className="bg-gray-800 p-6 rounded-xl shadow-lg mb-8">
          <div className="flex flex-col md:flex-row gap-4">
            <Input
              type="text"
              placeholder="输入关键词搜索数据..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="flex-grow bg-gray-700 border-gray-600"
            />
            <Button 
              onClick={handleSearch} 
              className="bg-blue-500 hover:bg-blue-600 flex items-center justify-center"
              disabled={isLoading}
            >
              <FaSearch className="mr-2" /> 
              {isLoading ? '搜索中...' : '搜索'}
            </Button>
            <Button 
              onClick={() => setShowFilters(!showFilters)} 
              className="bg-gray-700 hover:bg-gray-600 flex items-center justify-center"
            >
              <FaFilter className="mr-2" /> 筛选
            </Button>
            <select
              value={chainSource}
              onChange={(e) => setChainSource(e.target.value)}
              className="bg-gray-700 border border-gray-600 rounded-md p-2"
            >
              <option value="all">所有区块链</option>
              <option value="ethereum">仅以太坊</option>
              <option value="fabric">仅Fabric</option>
            </select>
          </div>
          
          {/* 高级筛选选项 */}
          {showFilters && (
            <div className="mt-4 p-4 bg-gray-700 rounded-lg grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">数据类型</label>
                <select
                  className="w-full p-2 rounded-md bg-gray-600 border-gray-500"
                  value={selectedType}
                  onChange={(e) => setSelectedType(e.target.value)}
                >
                  <option value="all">所有类型</option>
                  <option value="电子病历">电子病历</option>
                  <option value="基因组数据">基因组数据</option>
                  <option value="影像数据">影像数据</option>
                  <option value="处方数据">处方数据</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-1">上传时间范围</label>
                <div className="flex gap-2">
                  <Input
                    type="date"
                    value={dateRange.start}
                    onChange={(e) => setDateRange({...dateRange, start: e.target.value})}
                    className="w-full bg-gray-600 border-gray-500"
                  />
                  <span className="self-center">至</span>
                  <Input
                    type="date"
                    value={dateRange.end}
                    onChange={(e) => setDateRange({...dateRange, end: e.target.value})}
                    className="w-full bg-gray-600 border-gray-500"
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-1">排序方式</label>
                <select
                  className="w-full p-2 rounded-md bg-gray-600 border-gray-500"
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value)}
                >
                  <option value="newest">最新上传</option>
                  <option value="oldest">最早上传</option>
                  <option value="name">文件名</option>
                </select>
              </div>
              
              <div className="md:col-span-3 flex justify-end">
                <Button 
                  onClick={resetFilters} 
                  className="bg-gray-600 hover:bg-gray-500 mr-2"
                >
                  重置筛选
                </Button>
                <Button 
                  onClick={handleSearch} 
                  className="bg-blue-500 hover:bg-blue-600"
                >
                  应用筛选
                </Button>
              </div>
            </div>
          )}
        </Card>
        
        {/* 搜索结果 */}
        <div>
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-xl font-bold">搜索结果</h3>
            <p className="text-gray-400">{searchResults.length} 条记录</p>
          </div>
          
          {isLoading ? (
            <div className="flex justify-center items-center h-64">
              <p className="text-xl">加载中...</p>
            </div>
          ) : searchResults.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {searchResults.map((item) => (
                <Card key={item.id} className="bg-gray-800 hover:bg-gray-750 transition-colors overflow-hidden">
                  <div className="p-4">
                    <div className="flex items-start mb-3">
                      <div className="p-3 bg-gray-700 rounded-full mr-3">
                        {getDataTypeIcon(item.dataType)}
                      </div>
                      <div>
                        <h4 className="font-bold truncate">{item.fileName}</h4>
                        <div className="flex items-center">
                          <p className="text-sm text-gray-400">{item.dataType} • {item.fileSize}</p>
                          <span className="ml-2 flex items-center text-xs px-2 py-1 rounded-full bg-gray-700">
                            {item.chain === 'ethereum' ? 
                              <><FaEthereum className="mr-1 text-blue-400" /> 以太坊</> : 
                              <><FaNetworkWired className="mr-1 text-green-400" /> Fabric</>}
                          </span>
                        </div>
                      </div>
                    </div>
                    
                    <p className="text-sm text-gray-300 mb-3 line-clamp-2">{item.description}</p>
                    
                    <div className="flex flex-wrap gap-1 mb-3">
                      {item.tags.map((tag, index) => (
                        <span key={index} className="text-xs bg-blue-900 text-blue-200 px-2 py-1 rounded-full">
                          {tag}
                        </span>
                      ))}
                    </div>
                    
                    <div className="flex justify-between items-center text-sm text-gray-400">
                      <span>{item.owner} • {item.hospital}</span>
                      <span>{new Date(item.uploadTime).toLocaleDateString()}</span>
                    </div>
                  </div>
                  
                  <div className="bg-gray-700 p-3 flex justify-between items-center">
                    <span className="text-sm">{item.uploadTime}</span>
                    <Link to={`/data-detail/${item.id}`} className="flex items-center text-blue-400 hover:text-blue-300">
                      <FaEye className="mr-1" /> 查看详情
                    </Link>
                  </div>
                </Card>
              ))}
            </div>
          ) : (
            <div className="bg-gray-800 rounded-xl p-8 text-center">
              <p className="text-xl mb-2">未找到匹配的数据</p>
              <p className="text-gray-400">请尝试调整搜索条件或筛选条件</p>
            </div>
          )}
        </div>
        
        {/* 查询指南 */}
        <Card className="bg-gray-800 p-6 rounded-xl shadow-lg mt-8">
          <h3 className="text-xl font-bold mb-4">数据查询指南</h3>
          <div className="text-gray-300 space-y-2">
            <p>1. 您可以使用关键词搜索数据的文件名、描述或标签。</p>
            <p>2. 使用高级筛选可以按数据类型、上传时间等条件精确查找。</p>
            <p>3. 点击"查看详情"可以查看数据的完整信息，包括预览、元数据和授权历史。</p>
            <p>4. 所有数据查询操作都会通过区块链记录，确保数据访问的透明性和可追溯性。</p>
            <p>5. 您只能查看您有权访问的数据，包括您上传的数据和他人授权给您的数据。</p>
          </div>
        </Card>
      </div>
    </div>
  );
}