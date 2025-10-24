"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function DashboardPage() {
  const router = useRouter();
  const [isProcessing, setIsProcessing] = useState(false);
  const [loginCaixa, setLoginCaixa] = useState("");
  const [senhaCaixa, setSenhaCaixa] = useState("");
  const [cpfConsulta, setCpfConsulta] = useState("");

  const handleConsulta = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsProcessing(true);

    // Simulando processamento
    setTimeout(() => {
      setIsProcessing(false);
      alert("Consulta processada com sucesso! PDFs gerados.");
      // Limpar campos
      setCpfConsulta("");
    }, 3000);
  };

  const handleLogout = () => {
    router.push("/");
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-indigo-600 rounded-lg flex items-center justify-center">
                <svg
                  className="w-6 h-6 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
              </div>
              <div>
                <h1 className="text-xl font-bold text-gray-900">
                  Automação Caixa
                </h1>
                <p className="text-sm text-gray-500">Dashboard</p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <button
                onClick={() => router.push("/historico")}
                className="px-4 py-2 text-gray-700 hover:text-indigo-600 font-medium transition-colors"
              >
                📊 Histórico
              </button>
              <button
                onClick={handleLogout}
                className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-medium transition-colors"
              >
                Sair
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Card de Boas-vindas */}
        <div className="bg-gradient-to-r from-indigo-500 to-purple-600 rounded-2xl p-8 text-white mb-8 shadow-lg">
          <h2 className="text-3xl font-bold mb-2">Bem-vindo de volta! 👋</h2>
          <p className="text-indigo-100">
            Faça sua consulta e gere os documentos automaticamente
          </p>
        </div>

        {/* Formulário Principal */}
        <div className="bg-white rounded-2xl shadow-lg p-8 mb-8">
          <div className="mb-6">
            <h3 className="text-2xl font-bold text-gray-900 mb-2">
              Nova Consulta
            </h3>
            <p className="text-gray-600">
              Preencha os dados para iniciar a automação
            </p>
          </div>

          <form onSubmit={handleConsulta} className="space-y-6">
            {/* Credenciais da Caixa */}
            <div className="border-b border-gray-200 pb-6">
              <h4 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
                <span className="bg-indigo-100 text-indigo-600 w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold">
                  1
                </span>
                Credenciais Portal Caixa
              </h4>

              <div className="grid md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Login Caixa
                  </label>
                  <input
                    type="text"
                    value={loginCaixa}
                    onChange={(e) => setLoginCaixa(e.target.value)}
                    required
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none"
                    placeholder="Seu login da Caixa"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Senha Caixa
                  </label>
                  <input
                    type="password"
                    value={senhaCaixa}
                    onChange={(e) => setSenhaCaixa(e.target.value)}
                    required
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none"
                    placeholder="••••••••"
                  />
                </div>
              </div>
            </div>

            {/* CPF para Consulta */}
            <div>
              <h4 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
                <span className="bg-indigo-100 text-indigo-600 w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold">
                  2
                </span>
                Dados para Consulta
              </h4>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  CPF do Cliente
                </label>
                <input
                  type="text"
                  value={cpfConsulta}
                  onChange={(e) => setCpfConsulta(e.target.value)}
                  required
                  maxLength={14}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all outline-none text-lg"
                  placeholder="000.000.000-00"
                />
                <p className="text-sm text-gray-500 mt-2">
                  Digite o CPF do cliente que deseja consultar
                </p>
              </div>
            </div>

            {/* Botão de Processar */}
            <div className="pt-4">
              <button
                type="submit"
                disabled={isProcessing}
                className="w-full bg-indigo-600 text-white py-4 px-6 rounded-lg font-semibold text-lg hover:bg-indigo-700 focus:ring-4 focus:ring-indigo-200 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-lg"
              >
                {isProcessing ? (
                  <span className="flex items-center justify-center">
                    <svg
                      className="animate-spin -ml-1 mr-3 h-6 w-6 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      ></circle>
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                    Processando... Isso pode levar até 2 minutos
                  </span>
                ) : (
                  <span className="flex items-center justify-center">
                    🚀 Iniciar Consulta e Gerar Documentos
                  </span>
                )}
              </button>
            </div>
          </form>
        </div>

        {/* Cards Informativos */}
        <div className="grid md:grid-cols-3 gap-6">
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <div className="text-3xl mb-2">⚡</div>
            <h4 className="font-semibold text-gray-900 mb-1">Rápido</h4>
            <p className="text-sm text-gray-600">
              Processamento em menos de 2 minutos
            </p>
          </div>

          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <div className="text-3xl mb-2">📄</div>
            <h4 className="font-semibold text-gray-900 mb-1">3 Documentos</h4>
            <p className="text-sm text-gray-600">Geração automática em PDF</p>
          </div>

          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <div className="text-3xl mb-2">🔒</div>
            <h4 className="font-semibold text-gray-900 mb-1">Seguro</h4>
            <p className="text-sm text-gray-600">Dados criptografados</p>
          </div>
        </div>
      </main>
    </div>
  );
}
