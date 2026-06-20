import React, { useEffect, useState, useCallback } from 'react';
import { userService } from '../services/userService';
import type { User } from '../services/userService';
import { Edit, Trash2, Plus, ChevronLeft, ChevronRight, Loader2 } from 'lucide-react';
import toast from 'react-hot-toast';

interface UserListProps {
  onEdit: (user: User) => void;
  onAdd: () => void;
  refreshTrigger: number;
}

const UserList: React.FC<UserListProps> = ({ onEdit, onAdd, refreshTrigger }) => {
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [loading, setLoading] = useState(false);

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    try {
      const response = await userService.getAll(page, limit);
      if (response.success) {
        setUsers(response.data.users);
        setTotal(response.data.total);
      }
    } catch (err) {
      toast.error('Failed to fetch users');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [page, limit]);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchUsers();
  }, [fetchUsers, page, limit, refreshTrigger]);

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        const response = await userService.delete(id);
        if (response.success) {
          toast.success('User deleted successfully');
          fetchUsers();
        }
      } catch {
        toast.error('Failed to delete user');
      }
    }
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
      <div className="p-6 border-b border-slate-200 flex justify-between items-center bg-slate-50/50">
        <div>
          <h2 className="text-xl font-semibold text-slate-800">Users</h2>
          <p className="text-sm text-slate-500 mt-1">Manage your organization's users</p>
        </div>
        <button
          onClick={onAdd}
          className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg transition-all shadow-sm active:scale-95"
        >
          <Plus size={20} />
          <span>Add User</span>
        </button>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-50 border-b border-slate-200 text-slate-500 text-sm font-medium">
              <th className="px-6 py-4">Name</th>
              <th className="px-6 py-4">Email</th>
              <th className="px-6 py-4">Phone</th>
              <th className="px-6 py-4">PAN / Aadhaar</th>
              <th className="px-6 py-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {loading ? (
              <tr>
                <td colSpan={5} className="px-6 py-12 text-center text-slate-400">
                  <div className="flex flex-col items-center gap-2">
                    <Loader2 className="animate-spin text-indigo-600" size={32} />
                    <span>Loading users...</span>
                  </div>
                </td>
              </tr>
            ) : users.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-12 text-center text-slate-400">
                  No users found
                </td>
              </tr>
            ) : (
              users.map((user) => (
                <tr key={user.id} className="hover:bg-slate-50/50 transition-colors group">
                  <td className="px-6 py-4">
                    <div className="font-medium text-slate-900">{user.name}</div>
                  </td>
                  <td className="px-6 py-4 text-slate-600 font-mono text-sm">{user.email}</td>
                  <td className="px-6 py-4 text-slate-600 font-mono text-sm">{user.primaryMobile}</td>
                  <td className="px-6 py-4">
                    <div className="text-xs space-y-1">
                      <div className="text-slate-500 uppercase tracking-wider">PAN: <span className="text-slate-700 font-mono">{user.pan}</span></div>
                      <div className="text-slate-500 uppercase tracking-wider">AAD: <span className="text-slate-700 font-mono">{user.aadhaar}</span></div>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button
                        onClick={() => onEdit(user)}
                        className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all"
                        title="Edit User"
                      >
                        <Edit size={18} />
                      </button>
                      <button
                        onClick={() => handleDelete(user.id)}
                        className="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                        title="Delete User"
                      >
                        <Trash2 size={18} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <div className="p-4 border-t border-slate-200 flex items-center justify-between bg-slate-50/50 text-sm text-slate-600">
        <div>
          Showing <span className="font-medium">{users.length}</span> of <span className="font-medium">{total}</span> users
        </div>
        <div className="flex gap-2">
          <button
            disabled={page === 1}
            onClick={() => setPage(page - 1)}
            className="p-2 border border-slate-300 rounded-lg hover:bg-white disabled:opacity-50 disabled:bg-slate-100 transition-all shadow-sm"
          >
            <ChevronLeft size={18} />
          </button>
          <div className="flex items-center px-4 font-medium">
            Page {page} of {totalPages || 1}
          </div>
          <button
            disabled={page >= totalPages}
            onClick={() => setPage(page + 1)}
            className="p-2 border border-slate-300 rounded-lg hover:bg-white disabled:opacity-50 disabled:bg-slate-100 transition-all shadow-sm"
          >
            <ChevronRight size={18} />
          </button>
        </div>
      </div>
    </div>
  );
};

export default UserList;
