import React, { useState } from 'react';
import { Button } from "../components/ui/button"; // 按钮组件
import { Input } from "../components/ui/input"; // 输入框组件
import { Card } from "../components/ui/card"; // 卡片组件

export default function DataPage() {
  const [file, setFile] = useState(null);
  const [dataType, setDataType] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [uploadedData, setUploadedData] = useState([]);

  // 处理文件上传
  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = () => {
    if (!file || !dataType) {
      alert("请选择文件并选择数据类型！");
      return;
    }

    // 模拟上传过程（将文件数据和类型发送给后端）
    console.log("Uploading file:", file);
    console.log("Data type:", dataType);

    // 假设上传成功后，我们将数据添加到 uploadedData 列表
    setUploadedData([
      ...uploadedData,
      {
        id: uploadedData.length + 1,
        fileName: file.name,
        dataType,
        uploadTime: new Date().toLocaleString(),
      }
    ]);

    // 清空表单
    setFile(null);
    setDataType('');
  };

  // 处理查询
  const handleSearch = () => {
    console.log("Searching for:", searchQuery);
    // 在实际开发中这里会发送请求到后端并返回查询结果
    // 目前为示范数据展示
    const results = uploadedData.filter(item =>
      item.fileName.includes(searchQuery) || item.dataType.includes(searchQuery)
    );
    return results;
  };

  return (
    <div className="py-8 px-4">
      {/* 数据上传部分 */}
      <section className="mb-16">
        <h2 className="text-3xl font-semibold text-center mb-4">数据上传</h2>
        <div className="flex justify-center gap-8">
          <input
            type="file"
            onChange={handleFileChange}
            className="px-4 py-2 border border-gray-300 rounded-md"
          />
          <select
            className="px-4 py-2 border border-gray-300 rounded-md"
            value={dataType}
            onChange={(e) => setDataType(e.target.value)}
          >
            <option value="">选择数据类型</option>
            <option value="电子病历">电子病历</option>
            <option value="基因组数据">基因组数据</option>
            <option value="影像数据">影像数据</option>
          </select>
          <Button onClick={handleUpload} className="bg-blue-500 hover:bg-blue-600 text-white">上传数据</Button>
        </div>
      </section>

      {/* 数据查询部分 */}
      <section>
        <h2 className="text-3xl font-semibold text-center mb-4">数据查询</h2>
        <div className="flex justify-center gap-4 mb-6">
          <Input
            type="text"
            placeholder="请输入查询条件"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-md"
          />
          <Button onClick={handleSearch} className="bg-green-500 hover:bg-green-600 text-white">查询</Button>
        </div>

        {/* 显示查询结果 */}
        <div className="flex justify-center gap-8">
          {handleSearch().map((item) => (
            <Card key={item.id} className="bg-gray-800 p-6 w-72">
              <h4 className="text-white font-bold">{item.fileName}</h4>
              <p className="text-gray-400">{item.dataType}</p>
              <p className="text-gray-400">{item.uploadTime}</p>
            </Card>
          ))}
        </div>
      </section>
    </div>
  );
}
