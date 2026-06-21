import React, { useState } from 'react';
import { userService } from '../services/userService';
import type { CreateUserRequest, User } from '../services/userService';
import { X, Save, Loader2 } from 'lucide-react';
import toast from 'react-hot-toast';

interface UserFormProps {
  user?: User;
  onClose: () => void;
  onSuccess: () => void;
}

const UserForm: React.FC<UserFormProps> = ({ user, onClose, onSuccess }) => {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<CreateUserRequest>({
    name: user?.name || '',
    email: user?.email || '',
    primaryMobile: user?.primaryMobile || '',
    secondaryMobile: user?.secondaryMobile || '',
    aadhaar: user?.aadhaar || '',
    pan: user?.pan || '',
    dateOfBirth: user?.dateOfBirth || '',
    placeOfBirth: user?.placeOfBirth || '',
    currentAddress: user?.currentAddress || '',
    permanentAddress: user?.permanentAddress || '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      if (user) {
        await userService.update(user.id, formData);
        toast.success('User updated successfully');
      } else {
        await userService.create(formData);
        toast.success('User created successfully');
      }
      onSuccess();
    } catch (error: unknown) {
      const message = (error as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Something went wrong';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4 z-50 animate-in fade-in duration-200">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-hidden flex flex-col animate-in zoom-in-95 duration-200">
        <div className="p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
          <div>
            <h2 className="text-xl font-bold text-slate-800">{user ? 'Edit User' : 'Add New User'}</h2>
            <p className="text-sm text-slate-500 mt-1">Please fill in the user details below</p>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-slate-200 rounded-full transition-colors text-slate-400 hover:text-slate-600">
            <X size={20} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="overflow-y-auto p-6 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">Full Name</label>
              <input
                required
                name="name"
                value={formData.name}
                onChange={handleChange}
                className="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all outline-none"
                placeholder="John Doe"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">Email Address</label>
              <input
                required
                type="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all outline-none"
                placeholder="john@example.com"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">Primary Mobile</label>
              <input
                required
                name="primaryMobile"
                value={formData.primaryMobile}
                onChange={handleChange}
                className="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all outline-none"
                placeholder="9876543210"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">Secondary Mobile (Optional)</label>
              <input
                name="secondaryMobile"
                value={formData.secondaryMobile}
                onChange={handleChange}
                className="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all outline-none"
                placeholder="9876543211"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">Aadhaar Number</label>
              <input
                disabled={!!user}
                required={!user}
                name="aadhaar"
                value={formData.aadhaar}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg transition-all outline-none ${
                    user 
                    ? 'bg-slate-100 border-slate-200 text-slate-500 cursor-not-allowed' 
                    : 'bg-slate-50 border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500'
                }`}
                placeholder="1234 5678 9012"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">PAN Number</label>
              <input
                disabled={!!user}
                required={!user}
                name="pan"
                value={formData.pan}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg transition-all outline-none ${
                    user 
                    ? 'bg-slate-100 border-slate-200 text-slate-500 cursor-not-allowed' 
                    : 'bg-slate-50 border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500'
                }`}
                placeholder="ABCDE1234F"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">Date of Birth</label>
              <input
                required
                type="date"
                name="dateOfBirth"
                value={formData.dateOfBirth}
                onChange={handleChange}
                className="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all outline-none"
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">Place of Birth</label>
              <input
                required
                name="placeOfBirth"
                value={formData.placeOfBirth}
                onChange={handleChange}
                className="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all outline-none"
                placeholder="Mumbai"
              />
            </div>
          </div>

          <div className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700">Current Address</label>
              <textarea
                required
                name="currentAddress"
                value={formData.currentAddress}
                onChange={handleChange}
                rows={3}
                className="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all outline-none resize-none"
                placeholder="123 Street Name..."
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-semibold text-slate-700 text-indigo-700 flex justify-between items-center">
                <span>Permanent Address</span>
                <button 
                  type="button" 
                  onClick={() => setFormData(prev => ({...prev, permanentAddress: prev.currentAddress}))}
                  className="text-xs font-normal hover:underline"
                >
                  Same as current
                </button>
              </label>
              <textarea
                required
                name="permanentAddress"
                value={formData.permanentAddress}
                onChange={handleChange}
                rows={3}
                className="w-full px-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all outline-none resize-none"
                placeholder="123 Street Name..."
              />
            </div>
          </div>
        </form>

        <div className="p-6 border-t border-slate-100 bg-slate-50 flex justify-end gap-3">
          <button
            type="button"
            onClick={onClose}
            className="px-6 py-2 text-slate-600 hover:text-slate-800 font-medium transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={loading}
            className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-2 rounded-lg font-medium transition-all shadow-md active:scale-95 disabled:bg-indigo-400"
          >
            {loading ? <Loader2 className="animate-spin" size={20} /> : <Save size={20} />}
            <span>{user ? 'Update' : 'Save'} User</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default UserForm;
