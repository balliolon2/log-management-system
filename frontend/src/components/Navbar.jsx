import { useNavigate } from "react-router-dom";
import { authService } from "../utils/auth";
import { LogOut, User, Shield } from "lucide-react";

function Navbar() {
  const navigate = useNavigate();
  const user = authService.getUser();

  const handleLogout = () => {
    authService.logout();
    navigate("/login");
  };

  if (!user) return null;

  return (
    <div className="bg-slate-800 border-b border-slate-700 px-6 py-4 flex justify-between items-center">
      <div className="flex items-center gap-3">
        <User className="text-blue-400 w-5 h-5" />
        <div>
          <div className="text-white font-semibold">{user.username}</div>
          <div className="text-xs text-slate-400 flex items-center gap-1">
            <Shield className="w-3 h-3" />
            {user.role === "admin" ? "Administrator" : "Viewer"} • Tenant:{" "}
            {user.tenant_id}
          </div>
        </div>
      </div>

      <button
        onClick={handleLogout}
        className="flex items-center gap-2 bg-red-500/20 hover:bg-red-500/30 text-red-400 px-4 py-2 rounded-lg transition duration-200"
      >
        <LogOut className="w-4 h-4" />
        Logout
      </button>
    </div>
  );
}

export default Navbar;
