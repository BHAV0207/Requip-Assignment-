import api from './api';

export interface User {
  id: string;
  name: string;
  email: string;
  primaryMobile: string;
  secondaryMobile?: string;
  aadhaar: string;
  pan: string;
  dateOfBirth: string;
  placeOfBirth: string;
  currentAddress: string;
  permanentAddress: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  primaryMobile: string;
  secondaryMobile?: string;
  aadhaar: string;
  pan: string;
  dateOfBirth: string;
  placeOfBirth: string;
  currentAddress: string;
  permanentAddress: string;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  primaryMobile?: string;
  secondaryMobile?: string;
  aadhaar?: string;
  pan?: string;
  dateOfBirth?: string;
  placeOfBirth?: string;
  currentAddress?: string;
  permanentAddress?: string;
}

export interface PaginatedResponse<T> {
  data: {
    users: T[];
    total: number;
    page: number;
    limit: number;
  };
  message: string;
  success: boolean;
}

export interface SingleResponse<T> {
  data: T;
  message: string;
  success: boolean;
}

export const userService = {
  getAll: async (page = 1, limit = 10) => {
    const response = await api.get<PaginatedResponse<User>>(`/users?page=${page}&limit=${limit}`);
    return response.data;
  },

  getById: async (id: string) => {
    const response = await api.get<SingleResponse<User>>(`/users/${id}`);
    return response.data;
  },

  create: async (req: CreateUserRequest) => {
    const response = await api.post<SingleResponse<User>>('/users', req);
    return response.data;
  },

  update: async (id: string, req: UpdateUserRequest) => {
    const response = await api.put<SingleResponse<User>>(`/users/${id}`, req);
    return response.data;
  },

  delete: async (id: string) => {
    const response = await api.delete<SingleResponse<null>>(`/users/${id}`);
    return response.data;
  },
};
