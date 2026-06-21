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

      <div className="p-6">
        {loading ? (
          <div className="flex flex-col items-center gap-2 py-12 text-slate-400">
            <Loader2 className="animate-spin text-indigo-600" size={32} />
            <span>Loading users...</span>
          </div>
        ) : users.length === 0 ? (
          <div className="py-12 text-center text-slate-400">No users found</div>
        ) : (
          <div className="space-y-4">
            {users.map((user) => (
              <div
                key={user.id}
                className="border border-slate-200 rounded-xl p-5 hover:border-indigo-200 hover:shadow-md transition-all bg-white"
              >
                {/* Header row: Name + Actions */}
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h3 className="text-lg font-semibold text-slate-900">{user.name}</h3>
                    <p className="text-sm text-slate-500 font-mono">{user.email}</p>
                  </div>
                  <div className="flex gap-1 shrink-0">
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
                </div>

                {/* Details grid */}
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-x-6 gap-y-3 text-sm">
                  <div>
                    <span className="text-slate-400 text-xs font-medium uppercase tracking-wider">Primary Mobile</span>
                    <p className="text-slate-700 font-mono mt-0.5">{user.primaryMobile}</p>
                  </div>
                  <div>
                    <span className="text-slate-400 text-xs font-medium uppercase tracking-wider">Secondary Mobile</span>
                    <p className="text-slate-700 font-mono mt-0.5">{user.secondaryMobile || '—'}</p>
                  </div>
                  <div>
                    <span className="text-slate-400 text-xs font-medium uppercase tracking-wider">Date of Birth</span>
                    <p className="text-slate-700 mt-0.5">{user.dateOfBirth || '—'}</p>
                  </div>
                  <div>
                    <span className="text-slate-400 text-xs font-medium uppercase tracking-wider">Place of Birth</span>
                    <p className="text-slate-700 mt-0.5">{user.placeOfBirth || '—'}</p>
                  </div>
                  <div>
                    <span className="text-slate-400 text-xs font-medium uppercase tracking-wider">PAN / Aadhaar</span>
                    <div className="mt-0.5 inline-flex items-center gap-1.5 bg-amber-50 text-amber-700 border border-amber-200 rounded-md px-2 py-1 text-xs font-medium">
                      <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                      Hidden for security
                    </div>
                  </div>
                  <div className="col-span-2 md:col-span-1 lg:col-span-1">
                    <span className="text-slate-400 text-xs font-medium uppercase tracking-wider">Current Address</span>
                    <p className="text-slate-700 mt-0.5">{user.currentAddress || '—'}</p>
                  </div>
                  <div className="col-span-2 md:col-span-1 lg:col-span-2">
                    <span className="text-slate-400 text-xs font-medium uppercase tracking-wider">Permanent Address</span>
                    <p className="text-slate-700 mt-0.5">{user.permanentAddress || '—'}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
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
