"use client";

import { useRouter } from "next/navigation";

// Dados fake para demonstração
const consultasMock = [
  {
    id: "1",
    cpf: "123.456.789-00",
    data: "23/10/2025 14:30",
    status: "Concluído",
    tempo: "1m 45s",
    documentos: 3,
  },
  {
    id: "2",
    cpf: "987.654.321-00",
    data: "23/10/2025 13:15",
    status: "Concluído",
    tempo: "1m 52s",
    documentos: 3,
  },
  {
    id: "3",
    cpf: "456.789.123-00",
    data: "23/10/2025 11:20",
    status: "Concluído",
    tempo: "1m 38s",
    documentos: 3,
  },
  {
    id: "4",
    cpf: "321.654.987-00",
    data: "22/10/2025 16:45",
    status: "Concluído",
    tempo: "2m 01s",
    documentos: 3,
  },
  {
    id: "5",
    cpf: "789.123.456-00",
    data: "22/10/2025 15:10",
    status: "Erro",
    tempo: "0m 30s",
    documentos: 0,
  },
  {
    id: "6",
    cpf: "654.321.789-00",
    data: "22/10/2025 10:05",
    status: "Concluído",
    tempo: "1m 55s",
    documentos: 3,
  },
];

export default function HistoricoPage() {
  const router = useRouter();

  const handleDownload = (id: string) => {
    alert(`Download dos documentos da consulta ${id} iniciado!`);
  };

  const getStatusColor = (status: string) => {
    if (status === "Concluído") return "bg-green-100 text-green-700";
    if (status === "Erro") return "bg-red-100 text-red-700";
    return "bg-yellow-100 text-yellow-700";
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
                <p className="text-sm text-gray-500">Histórico de Consultas</p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <button
                onClick={() => router.push("/dashboard")}
                className="px-4 py-2 text-gray-700 hover:text-indigo-600 font-medium transition-colors"
              >
                ← Voltar ao Dashboard
              </button>
              <button
                onClick={() => router.push("/")}
                className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-medium transition-colors"
              >
                Sair
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Stats Cards */}
        <div className="grid md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Total Consultas</p>
                <p className="text-3xl font-bold text-gray-900">
                  {consultasMock.length}
                </p>
              </div>
              <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                <span className="text-2xl">📊</span>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Concluídas</p>
                <p className="text-3xl font-bold text-green-600">
                  {consultasMock.filter((c) => c.status === "Concluído").length}
                </p>
              </div>
              <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
                <span className="text-2xl">✅</span>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Com Erro</p>
                <p className="text-3xl font-bold text-red-600">
                  {consultasMock.filter((c) => c.status === "Erro").length}
                </p>
              </div>
              <div className="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center">
                <span className="text-2xl">❌</span>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Tempo Médio</p>
                <p className="text-3xl font-bold text-indigo-600">1m 48s</p>
              </div>
              <div className="w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center">
                <span className="text-2xl">⏱️</span>
              </div>
            </div>
          </div>
        </div>

        {/* Tabela */}
        <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-xl font-bold text-gray-900">
              Consultas Recentes
            </h2>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    ID
                  </th>
                  <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    CPF Consultado
                  </th>
                  <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    Data/Hora
                  </th>
                  <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    Status
                  </th>
                  <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    Tempo
                  </th>
                  <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    Documentos
                  </th>
                  <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    Ações
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {consultasMock.map((consulta) => (
                  <tr
                    key={consulta.id}
                    className="hover:bg-gray-50 transition-colors"
                  >
                    <td className="px-6 py-4 text-sm text-gray-900 font-medium">
                      #{consulta.id}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">
                      {consulta.cpf}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {consulta.data}
                    </td>
                    <td className="px-6 py-4">
                      <span
                        className={`inline-flex px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(
                          consulta.status
                        )}`}
                      >
                        {consulta.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {consulta.tempo}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">
                      {consulta.documentos > 0 ? (
                        <span className="flex items-center gap-1">
                          📄 {consulta.documentos} PDFs
                        </span>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      {consulta.status === "Concluído" ? (
                        <button
                          onClick={() => handleDownload(consulta.id)}
                          className="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-sm font-medium transition-colors"
                        >
                          ⬇️ Baixar
                        </button>
                      ) : (
                        <button
                          disabled
                          className="px-3 py-1.5 bg-gray-100 text-gray-400 rounded-lg text-sm font-medium cursor-not-allowed"
                        >
                          Indisponível
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Paginação */}
          <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
            <p className="text-sm text-gray-600">
              Mostrando <strong>1-6</strong> de <strong>6</strong> consultas
            </p>
            <div className="flex gap-2">
              <button
                disabled
                className="px-3 py-1.5 bg-gray-100 text-gray-400 rounded-lg text-sm font-medium cursor-not-allowed"
              >
                ← Anterior
              </button>
              <button
                disabled
                className="px-3 py-1.5 bg-gray-100 text-gray-400 rounded-lg text-sm font-medium cursor-not-allowed"
              >
                Próxima →
              </button>
            </div>
          </div>
        </div>

        {/* Botão Nova Consulta */}
        <div className="mt-8 text-center">
          <button
            onClick={() => router.push("/dashboard")}
            className="px-8 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-semibold text-lg transition-colors shadow-lg"
          >
            ➕ Nova Consulta
          </button>
        </div>
      </main>
    </div>
  );
}
