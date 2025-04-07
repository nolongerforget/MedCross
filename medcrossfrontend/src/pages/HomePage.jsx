import { useState, useEffect } from "react";
import { Button } from "../components/ui/button";
import { Card, CardContent } from "../components/ui/card";
import { FaNetworkWired, FaFileMedical, FaSearch, FaChartBar, FaUserMd, FaHospital, FaLock, FaArrowRight } from "react-icons/fa";
import { FaEthereum } from "react-icons/fa6";
// import { SiHyperledger } from "react-icons/si";
import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { dataAPI } from "../services/api";

export default function HomePage() {
  // 统计数据状态
  const [statistics, setStatistics] = useState({
    totalRecords: 0,
    ethereumRecords: 0,
    fabricRecords: 0,
    dataTypes: {},
    monthlyGrowth: []
  });
  
  const [isLoading, setIsLoading] = useState(true);

  // 获取统计数据
  useEffect(() => {
    const fetchStatistics = async () => {
      try {
        const data = await dataAPI.getStatistics();
        setStatistics(data);
      } catch (error) {
        console.error('获取统计数据失败:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchStatistics();
  }, []);
  return (
    <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 text-white flex flex-col">
      {/* Navbar */}
      <nav className="flex justify-between items-center p-6 bg-opacity-80 bg-gray-800 shadow-lg w-full">
        <h1 className="text-2xl font-bold">MedCross</h1>
        <div className="space-x-6">
          <Link to="/" className="hover:text-blue-400">首页</Link>
          <Link to="/data-upload" className="hover:text-blue-400">数据上传</Link>
          <Link to="/data-query" className="hover:text-blue-400">数据查询</Link>
          <Link to="/profile" className="hover:text-blue-400">个人中心</Link>
          <Link to="/auth" className="bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded-md font-medium">登录/注册</Link>
        </div>
      </nav>
      
      {/* Hero Section */}
      <header className="text-center py-16 w-full">
        <motion.h2 initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 1 }}
          className="text-5xl font-extrabold text-blue-400">
          医疗数据跨链查询平台
        </motion.h2>
        <p className="mt-4 text-lg text-gray-300 max-w-3xl mx-auto">MedCross是一个安全、高效的医疗数据跨链共享平台，连接以太坊和Hyperledger Fabric区块链，实现医疗数据的安全共享与高效查询</p>
        <div className="flex justify-center gap-4 mt-6">
          <Link to="/data-query">
            <Button className="px-6 py-3 bg-blue-500 rounded-xl hover:bg-blue-600 flex items-center gap-2">
              立即查询 <FaSearch />
            </Button>
          </Link>
          <Link to="/auth">
            <Button className="px-6 py-3 bg-gray-700 rounded-xl hover:bg-gray-600 flex items-center gap-2">
              注册账号 <FaArrowRight />
            </Button>
          </Link>
        </div>
      </header>
      
      {/* 统计数据 */}
      <section className="py-8 w-full bg-gray-800 bg-opacity-50">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.1 }}
              className="bg-gray-700 rounded-lg p-6 shadow-lg"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-400 text-sm">总数据量</p>
                  <h3 className="text-3xl font-bold text-white">{statistics.totalRecords}</h3>
                </div>
                <FaChartBar className="text-blue-400 text-3xl" />
              </div>
              <div className="mt-4 text-sm text-gray-400">跨链医疗数据记录总数</div>
            </motion.div>
            
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.2 }}
              className="bg-gray-700 rounded-lg p-6 shadow-lg"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-400 text-sm">以太坊数据</p>
                  <h3 className="text-3xl font-bold text-white">{statistics.ethereumRecords}</h3>
                </div>
                <FaEthereum className="text-blue-400 text-3xl" />
              </div>
              <div className="mt-4 text-sm text-gray-400">存储在以太坊区块链上的记录</div>
            </motion.div>
            
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.3 }}
              className="bg-gray-700 rounded-lg p-6 shadow-lg"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-400 text-sm">Fabric数据</p>
                  <h3 className="text-3xl font-bold text-white">{statistics.fabricRecords}</h3>
                </div>
                {/* <SiHyperledger className="text-blue-400 text-3xl" /> */}
                <FaNetworkWired className="text-blue-400 text-3xl" />
              </div>
              <div className="mt-4 text-sm text-gray-400">存储在Hyperledger Fabric上的记录</div>
            </motion.div>
            
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.4 }}
              className="bg-gray-700 rounded-lg p-6 shadow-lg"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-400 text-sm">数据类型</p>
                  <h3 className="text-3xl font-bold text-white">{Object.keys(statistics.dataTypes).length}</h3>
                </div>
                <FaFileMedical className="text-blue-400 text-3xl" />
              </div>
              <div className="mt-4 text-sm text-gray-400">支持的医疗数据类型</div>
            </motion.div>
          </div>
        </div>
      </section>
      
      {/* Features */}
      <section id="features" className="py-16 text-center w-full">
        <h3 className="text-3xl font-bold text-blue-300">平台特点</h3>
        <div className="mt-8 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-6xl mx-auto px-4">
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            viewport={{ once: true }}
          >
            <Card className="bg-gray-800 p-6 h-full">
              <FaNetworkWired className="text-blue-400 text-4xl mx-auto" />
              <CardContent className="mt-4">
                <h4 className="text-xl font-bold mb-2">多链数据统一查询</h4>
                <p className="text-gray-400">同时连接以太坊和Hyperledger Fabric区块链网络，提供统一的查询接口，无需关心底层技术差异</p>
              </CardContent>
            </Card>
          </motion.div>
          
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.1 }}
            viewport={{ once: true }}
          >
            <Card className="bg-gray-800 p-6 h-full">
              <FaSearch className="text-blue-400 text-4xl mx-auto" />
              <CardContent className="mt-4">
                <h4 className="text-xl font-bold mb-2">跨链数据检索</h4>
                <p className="text-gray-400">强大的关键词和元数据搜索功能，支持多维度筛选，快速定位所需医疗数据</p>
              </CardContent>
            </Card>
          </motion.div>
          
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            viewport={{ once: true }}
          >
            <Card className="bg-gray-800 p-6 h-full">
              <FaFileMedical className="text-blue-400 text-4xl mx-auto" />
              <CardContent className="mt-4">
                <h4 className="text-xl font-bold mb-2">医疗数据上传与存储</h4>
                <p className="text-gray-400">支持多种医疗数据类型的上传和存储，包括影像数据、电子病历、基因组数据等</p>
              </CardContent>
            </Card>
          </motion.div>
          
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.3 }}
            viewport={{ once: true }}
          >
            <Card className="bg-gray-800 p-6 h-full">
              <FaLock className="text-blue-400 text-4xl mx-auto" />
              <CardContent className="mt-4">
                <h4 className="text-xl font-bold mb-2">安全可信的数据共享</h4>
                <p className="text-gray-400">基于区块链技术的不可篡改特性，确保医疗数据的真实性和完整性，支持数据溯源</p>
              </CardContent>
            </Card>
          </motion.div>
          
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.4 }}
            viewport={{ once: true }}
          >
            <Card className="bg-gray-800 p-6 h-full">
              <FaUserMd className="text-blue-400 text-4xl mx-auto" />
              <CardContent className="mt-4">
                <h4 className="text-xl font-bold mb-2">多角色支持</h4>
                <p className="text-gray-400">针对医生、研究人员、医院管理员和患者等不同角色，提供差异化的功能和权限管理</p>
              </CardContent>
            </Card>
          </motion.div>
          
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.5 }}
            viewport={{ once: true }}
          >
            <Card className="bg-gray-800 p-6 h-full">
              <FaHospital className="text-blue-400 text-4xl mx-auto" />
              <CardContent className="mt-4">
                <h4 className="text-xl font-bold mb-2">机构间协作</h4>
                <p className="text-gray-400">促进医疗机构之间的数据共享和协作，提高医疗资源利用效率，推动精准医疗发展</p>
              </CardContent>
            </Card>
          </motion.div>
        </div>
      </section>
      
      {/* 使用流程 */}
      <section id="workflow" className="py-16 w-full bg-gray-800 bg-opacity-50">
        <div className="container mx-auto px-4">
          <h3 className="text-3xl font-bold text-blue-300 text-center mb-12">使用流程</h3>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <motion.div 
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5 }}
              viewport={{ once: true }}
              className="flex flex-col items-center text-center"
            >
              <div className="bg-blue-500 rounded-full w-16 h-16 flex items-center justify-center mb-4">
                <span className="text-2xl font-bold">1</span>
              </div>
              <h4 className="text-xl font-bold mb-2">注册账号</h4>
              <p className="text-gray-400">创建MedCross账号，系统将自动为您生成以太坊钱包和Fabric身份</p>
              <Link to="/auth" className="mt-4 text-blue-400 hover:underline flex items-center gap-1">
                立即注册 <FaArrowRight size={12} />
              </Link>
            </motion.div>
            
            <motion.div 
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.2 }}
              viewport={{ once: true }}
              className="flex flex-col items-center text-center"
            >
              <div className="bg-blue-500 rounded-full w-16 h-16 flex items-center justify-center mb-4">
                <span className="text-2xl font-bold">2</span>
              </div>
              <h4 className="text-xl font-bold mb-2">上传数据</h4>
              <p className="text-gray-400">选择目标区块链，上传您的医疗数据，添加关键词和元数据以便后续检索</p>
              <Link to="/data-upload" className="mt-4 text-blue-400 hover:underline flex items-center gap-1">
                了解更多 <FaArrowRight size={12} />
              </Link>
            </motion.div>
            
            <motion.div 
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.4 }}
              viewport={{ once: true }}
              className="flex flex-col items-center text-center"
            >
              <div className="bg-blue-500 rounded-full w-16 h-16 flex items-center justify-center mb-4">
                <span className="text-2xl font-bold">3</span>
              </div>
              <h4 className="text-xl font-bold mb-2">跨链查询</h4>
              <p className="text-gray-400">使用强大的查询功能，跨链检索所需的医疗数据，支持多种筛选条件</p>
              <Link to="/data-query" className="mt-4 text-blue-400 hover:underline flex items-center gap-1">
                开始查询 <FaArrowRight size={12} />
              </Link>
            </motion.div>
          </div>
        </div>
      </section>      {/* 号召性用语 */}
      <section className="py-16 w-full">
        <div className="container mx-auto px-4 text-center">
          <motion.div 
            initial={{ opacity: 0, scale: 0.9 }}
            whileInView={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5 }}
            viewport={{ once: true }}
            className="bg-gradient-to-r from-blue-600 to-blue-800 rounded-xl p-10 max-w-4xl mx-auto shadow-2xl"
          >
            <h3 className="text-3xl font-bold mb-4">开始使用MedCross</h3>
            <p className="text-lg mb-6 max-w-2xl mx-auto">加入我们的医疗数据跨链共享平台，体验安全、高效的医疗数据管理与共享</p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link to="/auth">
                <Button className="bg-white text-blue-700 hover:bg-gray-100 px-8 py-3 text-lg font-medium">
                  注册账号
                </Button>
              </Link>
              <Link to="/data-query">
                <Button className="bg-transparent border-2 border-white hover:bg-blue-700 px-8 py-3 text-lg font-medium">
                  立即查询
                </Button>
              </Link>
            </div>
          </motion.div>
        </div>
      </section>
      
      {/* Footer */}
      <footer className="py-10 bg-gray-800 w-full">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
            <div>
              <h4 className="text-xl font-bold text-white mb-4">MedCross</h4>
              <p className="text-gray-400">医疗数据跨链共享平台，连接以太坊和Hyperledger Fabric区块链网络</p>
            </div>
            <div>
              <h4 className="text-lg font-bold text-white mb-4">功能</h4>
              <ul className="space-y-2 text-gray-400">
                <li><Link to="/data-upload" className="hover:text-blue-400">数据上传</Link></li>
                <li><Link to="/data-query" className="hover:text-blue-400">跨链查询</Link></li>
                <li><Link to="/profile" className="hover:text-blue-400">个人中心</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="text-lg font-bold text-white mb-4">资源</h4>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-blue-400">API文档</a></li>
                <li><a href="#" className="hover:text-blue-400">开发指南</a></li>
                <li><a href="#" className="hover:text-blue-400">常见问题</a></li>
              </ul>
            </div>
            <div>
              <h4 className="text-lg font-bold text-white mb-4">联系我们</h4>
              <p className="text-gray-400">有任何问题或建议，请随时联系我们</p>
              <p className="text-gray-400 mt-2">邮箱: contact@medcross.io</p>
            </div>
          </div>
          <div className="border-t border-gray-700 pt-6 text-center">
            <p className="text-gray-400">© 2025 MedCross. All Rights Reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
