"use client";

import { onAuthStateChanged, type User } from "firebase/auth";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";
import { AuthContext } from "@/shared/auth";
import { auth } from "@/shared/config";

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(auth.currentUser);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setUser(user);
      setLoading(false);
    });
    return () => unsubscribe();
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading }}>
      {children}
    </AuthContext.Provider>
  );
};
