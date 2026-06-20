import { useState } from 'react';
import UserList from './components/UserList';
import UserForm from './components/UserForm';
import type { User } from './services/userService';
import { Toaster } from 'react-hot-toast';
import { Users } from 'lucide-react';

function App() {
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | undefined>();
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleAdd = () => {
    setSelectedUser(undefined);
    setIsFormOpen(true);
  };

  const handleEdit = (user: User) => {
    setSelectedUser(user);
    setIsFormOpen(true);
  };

  const handleSuccess = () => {
    setIsFormOpen(false);
    setRefreshTrigger((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen bg-slate-50">
      <Toaster position="top-right" />
      
      {/* Navigation / Header */}
      <nav className="bg-white border-b border-slate-200 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <div className="flex items-center gap-3">
              <div className="bg-indigo-600 p-2 rounded-lg text-white">
                <Users size={24} />
              </div>
              <h1 className="text-xl font-bold text-slate-900 tracking-tight">
                Requip <span className="text-slate-500 font-normal">| User Manager</span>
              </h1>
            </div>
            <div className="hidden sm:flex items-center gap-4 text-sm font-medium text-slate-500">
              <span className="hover:text-indigo-600 cursor-pointer transition-colors">Dashboard</span>
              <span className="text-indigo-600 border-b-2 border-indigo-600 py-5">Users</span>
              <span className="hover:text-indigo-600 cursor-pointer transition-colors">Reports</span>
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h2 className="text-3xl font-extrabold text-slate-900">User Management</h2>
          <p className="text-slate-500 mt-2">View, edit, and manage your organization's user database in one place.</p>
        </div>

        <UserList 
          onEdit={handleEdit} 
          onAdd={handleAdd} 
          refreshTrigger={refreshTrigger} 
        />
      </main>

      {/* Form Modal */}
      {isFormOpen && (
        <UserForm 
          key={selectedUser?.id || 'new'}
          user={selectedUser} 
          onClose={() => setIsFormOpen(false)} 
          onSuccess={handleSuccess} 
        />
      )}
    </div>
  );
}

export default App;
