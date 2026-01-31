import { useState, useEffect } from "react";
import axios from "axios";
import { AlertTriangle, Clock, Shield, CheckCircle } from "lucide-react";

function AlertTable() {
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);

  // ฟังก์ชันดึง alerts จาก Backend
  const fetchAlerts = async () => {
    try {
      setLoading(true);
      const res = await axios.get(
        `/api/alerts?tenant=demo&limit=50`,
      );
      setAlerts(res.data || []);
    } catch (err) {
      console.error("Error fetching alerts:", err);
    } finally {
      setLoading(false);
    }
  };

  // ดึงข้อมูลครั้งแรก และ Auto-refresh ทุก 10 วินาที
  useEffect(() => {
    fetchAlerts();
    const interval = setInterval(() => {
      fetchAlerts();
    }, 10000); // Refresh ทุก 10 วิ
    return () => clearInterval(interval);
  }, []);

  // ฟังก์ชันแสดงสี Severity
  const getSeverityColor = (severity) => {
    if (severity >= 8) return "bg-red-500/20 text-red-400 border-red-500";
    if (severity >= 5)
      return "bg-orange-500/20 text-orange-400 border-orange-500";
    return "bg-yellow-500/20 text-yellow-400 border-yellow-500";
  };

  // ฟังก์ชันแสดง Icon ตาม Severity
  const getSeverityIcon = (severity) => {
    if (severity >= 8) return <AlertTriangle className="w-4 h-4" />;
    if (severity >= 5) return <Shield className="w-4 h-4" />;
    return <Clock className="w-4 h-4" />;
  };

  return (
    <div className="bg-slate-800 rounded-xl overflow-hidden shadow-xl border border-slate-700">
      {/* Header */}
      <div className="bg-slate-700 p-4 flex justify-between items-center">
        <h2 className="text-xl font-bold flex items-center gap-2">
          <AlertTriangle className="text-red-400" />
          Security Alerts
        </h2>
        <div className="text-sm text-slate-400">
          Total: <span className="font-bold text-white">{alerts.length}</span>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full text-left">
          <thead className="bg-slate-700 text-slate-300 text-sm">
            <tr>
              <th className="p-4">Severity</th>
              <th className="p-4">Rule Name</th>
              <th className="p-4">Details</th>
              <th className="p-4">Triggered At</th>
              <th className="p-4">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700">
            {loading ? (
              <tr>
                <td colSpan="5" className="p-8 text-center text-slate-500">
                  Loading alerts...
                </td>
              </tr>
            ) : alerts.length === 0 ? (
              <tr>
                <td colSpan="5" className="p-8 text-center text-slate-500">
                  <CheckCircle className="w-12 h-12 mx-auto mb-2 text-green-500" />
                  No alerts detected
                </td>
              </tr>
            ) : (
              alerts.map((alert, index) => (
                <tr
                  key={index}
                  className="hover:bg-slate-700/50 transition duration-150"
                >
                  {/* Severity */}
                  <td className="p-4">
                    <div
                      className={`flex items-center gap-2 px-3 py-1 rounded-lg border-l-4 ${getSeverityColor(
                        alert.severity,
                      )}`}
                    >
                      {getSeverityIcon(alert.severity)}
                      <span className="font-bold">{alert.severity}</span>
                    </div>
                  </td>

                  {/* Rule Name */}
                  <td className="p-4">
                    <div className="font-bold text-white">
                      {alert.rule_name}
                    </div>
                  </td>

                  {/* Details */}
                  <td className="p-4">
                    <div className="text-sm text-slate-300">
                      {alert.details?.message || "N/A"}
                    </div>
                    <div className="text-xs text-slate-500 mt-1">
                      {alert.details?.condition_field}=
                      {alert.details?.condition_value} ({alert.details?.count}/
                      {alert.details?.threshold})
                    </div>
                  </td>

                  {/* Triggered At */}
                  <td className="p-4">
                    <div className="text-sm text-slate-400 font-mono">
                      {new Date(alert.triggered_at).toLocaleString()}
                    </div>
                  </td>

                  {/* Status */}
                  <td className="p-4">
                    {alert.acknowledged ? (
                      <span className="bg-green-500/20 text-green-400 px-3 py-1 rounded text-xs font-bold">
                        Acknowledged
                      </span>
                    ) : (
                      <span className="bg-red-500/20 text-red-400 px-3 py-1 rounded text-xs font-bold animate-pulse">
                        Active
                      </span>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default AlertTable;
