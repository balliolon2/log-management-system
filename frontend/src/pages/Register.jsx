import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import axios from "axios";
import { UserPlus, Lock, User, Building, Shield } from "lucide-react";

function Register() {
    const [formData, setFormData] = useState({
        username: "",
        password: "",
        confirmPassword: "",
        role: "viewer",
        tenant_id: "default_tenant",
    });
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");

        if (formData.password !== formData.confirmPassword) {
            setError("Passwords do not match");
            return;
        }

        setLoading(true);

        try {
            // Backend expects: { username, password, role, tenant_id }
            const payload = {
                username: formData.username,
                password: formData.password,
                role: formData.role,
                tenant_id: formData.tenant_id
            };

            await axios.post("/api/register", payload);

            // Success - Redirect to login
            navigate("/login");
        } catch (err) {
            console.error("Register error:", err);
            setError(
                err.response?.data?.message || err.response?.data || "Registration failed. Please try again."
            );
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
                        <UserPlus className="w-16 h-16 text-blue-400" />
                    </div>
                    <h1 className="text-3xl font-bold text-white mb-2">Create Account</h1>
                    <p className="text-slate-400">Join Log Sentinel</p>
                </div>

                {/* Register Form */}
                <div className="bg-slate-800 rounded-xl shadow-2xl p-8 border border-slate-700">
                    <h2 className="text-2xl font-bold text-white mb-6">Sign Up</h2>

                    {error && (
                        <div className="bg-red-500/20 border border-red-500 text-red-400 px-4 py-3 rounded-lg mb-6">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-4">
                        {/* Username */}
                        <div>
                            <label className="block text-slate-300 text-sm font-medium mb-1">
                                Username
                            </label>
                            <div className="relative">
                                <User className="absolute left-3 top-3 text-slate-500 w-5 h-5" />
                                <input
                                    type="text"
                                    name="username"
                                    className="w-full bg-slate-700 border border-slate-600 text-white px-4 py-3 pl-12 rounded-lg focus:outline-none focus:border-blue-500 transition"
                                    placeholder="Choose a username"
                                    value={formData.username}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        {/* Password */}
                        <div>
                            <label className="block text-slate-300 text-sm font-medium mb-1">
                                Password
                            </label>
                            <div className="relative">
                                <Lock className="absolute left-3 top-3 text-slate-500 w-5 h-5" />
                                <input
                                    type="password"
                                    name="password"
                                    className="w-full bg-slate-700 border border-slate-600 text-white px-4 py-3 pl-12 rounded-lg focus:outline-none focus:border-blue-500 transition"
                                    placeholder="Create a password"
                                    value={formData.password}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        {/* Confirm Password */}
                        <div>
                            <label className="block text-slate-300 text-sm font-medium mb-1">
                                Confirm Password
                            </label>
                            <div className="relative">
                                <Lock className="absolute left-3 top-3 text-slate-500 w-5 h-5" />
                                <input
                                    type="password"
                                    name="confirmPassword"
                                    className="w-full bg-slate-700 border border-slate-600 text-white px-4 py-3 pl-12 rounded-lg focus:outline-none focus:border-blue-500 transition"
                                    placeholder="Confirm password"
                                    value={formData.confirmPassword}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        {/* Tenant ID */}
                        <div>
                            <label className="block text-slate-300 text-sm font-medium mb-1">
                                Organization ID (Tenant)
                            </label>
                            <div className="relative">
                                <Building className="absolute left-3 top-3 text-slate-500 w-5 h-5" />
                                <input
                                    type="text"
                                    name="tenant_id"
                                    className="w-full bg-slate-700 border border-slate-600 text-white px-4 py-3 pl-12 rounded-lg focus:outline-none focus:border-blue-500 transition"
                                    placeholder="e.g. default_tenant"
                                    value={formData.tenant_id}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        {/* Role */}
                        <div>
                            <label className="block text-slate-300 text-sm font-medium mb-1">
                                Role
                            </label>
                            <div className="relative">
                                <Shield className="absolute left-3 top-3 text-slate-500 w-5 h-5" />
                                <select
                                    name="role"
                                    className="w-full bg-slate-700 border border-slate-600 text-white px-4 py-3 pl-12 rounded-lg focus:outline-none focus:border-blue-500 transition appearance-none"
                                    value={formData.role}
                                    onChange={handleChange}
                                >
                                    <option value="viewer">Viewer</option>
                                    <option value="admin">Admin</option>
                                </select>
                                <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-4 text-slate-500">
                                    <svg className="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" /></svg>
                                </div>
                            </div>
                        </div>

                        {/* Submit Button */}
                        <button
                            type="submit"
                            disabled={loading}
                            className="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-3 px-4 rounded-lg transition duration-200 mt-4 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {loading ? "Creating Account..." : "Create Account"}
                        </button>
                    </form>

                    {/* Login Link */}
                    <div className="mt-6 text-center text-slate-400">
                        Already have an account?{" "}
                        <Link to="/login" className="text-blue-400 hover:text-blue-300 font-medium">
                            Sign In
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Register;
