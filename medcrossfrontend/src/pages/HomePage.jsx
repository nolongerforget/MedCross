import { Button } from "../components/ui/button";
import { Card, CardContent } from "../components/ui/card";
import { FaShieldAlt, FaNetworkWired, FaFileMedical } from "react-icons/fa";
import { motion } from "framer-motion";
import { Link } from "react-router-dom";

export default function HomePage() {
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
          <Link to="/blockchain-records" className="hover:text-blue-400">区块链记录</Link>
          <Link to="/auth" className="bg-blue-500 hover:bg-blue-600 px-3 py-1 rounded-md">登录/注册</Link>
        </div>
      </nav>
      
      {/* Hero Section */}
      <header className="text-center py-20 flex-grow w-full">
        <motion.h2 initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 1 }}
          className="text-5xl font-extrabold text-blue-400">
          医疗数据跨链共享新时代
        </motion.h2>
        <p className="mt-4 text-lg text-gray-300">安全、高效、合规的数据协作平台</p>
        <Button className="mt-6 px-6 py-3 bg-blue-500 rounded-xl hover:bg-blue-600">
          立即体验
        </Button>
      </header>
      
      {/* Features */}
      <section id="features" className="py-16 text-center w-full">
        <h3 className="text-3xl font-bold text-blue-300">平台特点</h3>
        <div className="mt-8 flex justify-center gap-8 w-full">
          <Card className="bg-gray-800 p-6 w-72">
            <FaShieldAlt className="text-blue-400 text-4xl mx-auto" />
            <CardContent className="mt-4">数据隐私与安全保障</CardContent>
          </Card>
          <Card className="bg-gray-800 p-6 w-72">
            <FaNetworkWired className="text-blue-400 text-4xl mx-auto" />
            <CardContent className="mt-4">区块链跨链数据互通</CardContent>
          </Card>
          <Card className="bg-gray-800 p-6 w-72">
            <FaFileMedical className="text-blue-400 text-4xl mx-auto" />
            <CardContent className="mt-4">医疗数据共享与溯源</CardContent>
          </Card>
        </div>
      </section>
      
      {/* Footer */}
      <footer className="text-center py-6 bg-gray-800 w-full">
        <p>© 2025 MedCross. All Rights Reserved.</p>
      </footer>
    </div>
  );
}
