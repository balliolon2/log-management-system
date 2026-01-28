import { useState, useEffect } from "react";
import axios from "axios";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Search, ShieldAlert, Activity, Server } from "lucide-react";
import AlertTable from "../components/AlertTable"; // แก้จาก "./components/AlertTable"
import Navbar from "../components/Navbar"; // เพิ่มบรรทัดนี้

function Dashboard() {
  // เปลี่ยนจาก App เป็น Dashboard
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");

  // ฟังก์ชันดึงข้อมูลจาก Backend
  const fetchLogs = async (query = "") => {
    try {
      setLoading(true);
      // ยิงไปที่ Backend Go ของเรา (Port 8080)
      const res = await axios.get(
        `http://localhost:8080/api/logs?q=${query}&limit=50`,
      );
      setLogs(res.data || []);
    } catch (err) {
      console.error("Error fetching logs:", err);
    } finally {
      setLoading(false);
    }
  };

  // ดึงข้อมูลครั้งแรก และ Auto-refresh
  useEffect(() => {
    fetchLogs();
    const interval = setInterval(() => {
      if (searchTerm === "") fetchLogs();
    }, 5000); // Refresh ทุก 5 วิ
    return () => clearInterval(interval);
  }, [searchTerm]);

  const handleSearch = (e) => {
    if (e.key === "Enter") fetchLogs(searchTerm);
  };

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100">
      {/* Navbar - เพิ่มบรรทัดนี้ */}
      <Navbar />

      <div className="p-8 font-sans">
        {/* Header */}
        <header className="mb-8 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-blue-400 flex items-center gap-2">
              <ShieldAlert className="w-8 h-8" />
              Log Sentinel{" "}
              <span className="text-xs text-slate-500 bg-slate-800 px-2 py-1 rounded">
                v4.0
              </span>
            </h1>
            <p className="text-slate-400 mt-1">
              Real-time Security Event Monitoring
            </p>
          </div>
          <div className="flex gap-4">
            <div className="bg-slate-800 p-4 rounded-lg flex items-center gap-3">
              <Server className="text-green-400" />
              <div>
                <div className="text-xs text-slate-400">System Status</div>
                <div className="font-bold text-green-400">Online</div>
              </div>
            </div>
          </div>
        </header>

        {/* Search Bar */}
        <div className="mb-6 relative">
          <input
            type="text"
            placeholder="Search logs..."
            className="w-full bg-slate-800 border border-slate-700 text-white px-4 py-3 pl-12 rounded-lg focus:outline-none focus:border-blue-500 transition"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onKeyDown={handleSearch}
          />
          <Search className="absolute left-4 top-3.5 text-slate-500 w-5 h-5" />
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-slate-800 p-6 rounded-xl border-l-4 border-blue-500">
            <h3 className="text-slate-400 text-sm mb-1">Total Logs</h3>
            <div className="text-3xl font-bold">{logs.length}</div>
          </div>
          <div className="bg-slate-800 p-6 rounded-xl border-l-4 border-red-500">
            <h3 className="text-slate-400 text-sm mb-1">Critical Events</h3>
            <div className="text-3xl font-bold">
              {logs.filter((l) => l.severity > 3).length}
            </div>
          </div>
          <div className="bg-slate-800 p-6 rounded-xl border-l-4 border-emerald-500">
            <h3 className="text-slate-400 text-sm mb-1">Active Sources</h3>
            <div className="text-3xl font-bold">
              {[...new Set(logs.map((l) => l.source))].length}
            </div>
          </div>
        </div>

        {/* Alert Table - เพิ่มส่วนนี้ */}
        <div className="mb-8">
          <AlertTable />
        </div>

        {/* Graph */}
        <div className="bg-slate-800 p-6 rounded-xl mb-8">
          <h2 className="text-xl font-bold mb-4 flex items-center gap-2">
            <Activity className="text-blue-400" /> Event Traffic
          </h2>
          <div className="h-64 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={logs.slice().reverse()}>
                <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
                <XAxis dataKey="@timestamp" hide />
                <YAxis stroke="#94a3b8" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "#1e293b",
                    borderColor: "#334155",
                    color: "#fff",
                  }}
                />
                <Line
                  type="monotone"
                  dataKey="severity"
                  stroke="#3b82f6"
                  strokeWidth={2}
                  dot={false}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Log Table */}
        <div className="bg-slate-800 rounded-xl overflow-hidden shadow-xl border border-slate-700">
          <table className="w-full text-left">
            <thead className="bg-slate-700 text-slate-300">
              <tr>
                <th className="p-4">Timestamp</th>
                <th className="p-4">Source</th>
                <th className="p-4">Message</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-700">
              {loading ? (
                <tr>
                  <td colSpan="3" className="p-8 text-center text-slate-500">
                    Loading...
                  </td>
                </tr>
              ) : logs.length === 0 ? (
                <tr>
                  <td colSpan="3" className="p-8 text-center text-slate-500">
                    No logs found
                  </td>
                </tr>
              ) : (
                logs.map((log, index) => (
                  <tr
                    key={index}
                    className="hover:bg-slate-700/50 transition duration-150"
                  >
                    <td className="p-4 text-sm text-slate-400 font-mono">
                      {new Date(log["@timestamp"]).toLocaleTimeString()}
                    </td>
                    <td className="p-4">
                      <span className="bg-blue-500/20 text-blue-400 px-2 py-1 rounded text-xs font-bold uppercase">
                        {log.source}
                      </span>
                    </td>
                    <td className="p-4 text-sm text-slate-300 font-mono truncate max-w-md">
                      {log.raw}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
