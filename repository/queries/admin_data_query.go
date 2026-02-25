package queries

// AdminDataQuery is the SELECT used for admin/internal tokens (is_super_token = true).
// Edit this file to customize the admin query — columns, JOINs, ORDER BY, etc. —
// without touching any business logic or repository code.
//
// Rules:
//   - Must SELECT exactly 27 columns in the same order as scanDataRow() in data_repository.go:
//     op cols (23): Terminal ID, Terminal Name, Priority, Mode, Initial Problem,
//                   Current Problem, P-Duration, Incident start datetime, Count,
//                   Status, Remarks, Balance, Condition, Tickets no, Tickets duration,
//                   Open time, Close time, Problem History, Mode History,
//                   DSP FLM, DSP SLM, Last Withdrawal, Export Name
//     mm cols  (4): FLM name, FLM, SLM, Net
//   - Do NOT include WHERE or ORDER BY here; the repository appends them for pagination.
//     If you want a fixed sort, add ORDER BY before the final comment line.
const AdminDataQuery = `
	SELECT
		op.[Terminal ID],
		op.[Terminal Name],
		op.[Priority],
		op.[Mode],
		op.[Initial Problem],
		op.[Current Problem],
		op.[P-Duration],
		op.[Incident start datetime],
		op.[Count],
		op.[Status],
		op.[Remarks],
		op.[Balance],
		op.[Condition],
		op.[Tickets no],
		op.[Tickets duration],
		op.[Open time],
		op.[Close time],
		op.[Problem History],
		op.[Mode History],
		op.[DSP FLM],
		op.[DSP SLM],
		op.[Last Withdrawal],
		op.[Export Name],
		mm.[FLM name],
		mm.[FLM],
		mm.[SLM],
		mm.[Net]
	FROM ticket_master.dbo.open_ticket op
	LEFT JOIN machine_master.dbo.machine mm
		ON op.[Terminal ID] = mm.[Terminal ID]
`
