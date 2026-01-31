import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import axios from "axios";
import { authService } from "../utils/auth";
import { ShieldAlert, Lock, User } from "lucide-react";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await axios.post("/api/login", {
        username,
        password,
      });

      const { token, username: userName, role, tenant_id } = response.data;

      // บันทึก token และข้อมูล user
      authService.login(token, { username: userName, role, tenant_id });

      // ไปหน้า Dashboard
      navigate("/");
    } catch (err) {
      console.error("Login error:", err);
      setError(err.response?.data || "Invalid credentials. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-900 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        {/* Logo & Title */}
        <div className="text-center mb-8">
          <div className="flex justify-center mb-4">
            <ShieldAlert className="w-16 h-16 text-blue-400" />
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">Log Sentinel</h1>
          <p className="text-slate-400">Security Event Monitoring System</p>
        </div>

        {/* Login Form */}
        <div className="bg-slate-800 rounded-xl shadow-2xl p-8 border border-slate-700">
          <h2 className="text-2xl font-bold text-white mb-6">Sign In</h2>

          {error && (
            <div className="bg-red-500/20 border border-red-500 text-red-400 px-4 py-3 rounded-lg mb-6">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Username */}
            <div>
              <label className="block text-slate-300 text-sm font-medium mb-2">
                Username
              </label>
              <div className="relative">
                <User className="absolute left-3 top-3 text-slate-500 w-5 h-5" />
                <input
                  type="text"
                  className="w-full bg-slate-700 border border-slate-600 text-white px-4 py-3 pl-12 rounded-lg focus:outline-none focus:border-blue-500 transition"
                  placeholder="Enter your username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </div>
            </div>

            {/* Password */}
            <div>
              <label className="block text-slate-300 text-sm font-medium mb-2">
                Password
              </label>
              <div className="relative">
                <Lock className="absolute left-3 top-3 text-slate-500 w-5 h-5" />
                <input
                  type="password"
                  className="w-full bg-slate-700 border border-slate-600 text-white px-4 py-3 pl-12 rounded-lg focus:outline-none focus:border-blue-500 transition"
                  placeholder="Enter your password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-3 px-4 rounded-lg transition duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? "Signing in..." : "Sign In"}
            </button>
          </form>

          {/* Register Link */}
          <div className="mt-6 text-center text-slate-400">
            Don't have an account?{" "}
            <Link to="/register" className="text-blue-400 hover:text-blue-300 font-medium">
              Create Account
            </Link>
          </div>

          {/* Demo Credentials */}
          <div className="mt-6 p-4 bg-slate-700/50 rounded-lg border border-slate-600">
            <p className="text-slate-400 text-xs mb-2 font-semibold">
              Demo Credentials:
            </p>
            <div className="text-xs text-slate-300 space-y-1">
              <p>
                <span className="text-blue-400">Admin:</span> admin / admin123
              </p>
              <p>
                <span className="text-green-400">Viewer:</span> viewer1 /
                viewer123
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Login;
