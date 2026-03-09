-- 000008_runbook-step-refactor.sql: Rename step types and migrate check→run.
--
-- Step type renames:
--   command → run
--   check   → run  (moves "check" field value into "command" field)
--   manual  → approval

-- Step 1: Rename "command" type to "run".
UPDATE ops_runbooks SET steps_json = REPLACE(
    steps_json,
    '"type":"command"',
    '"type":"run"'
) WHERE steps_json LIKE '%"type":"command"%';

-- Step 2: Rename "manual" type to "approval".
UPDATE ops_runbooks SET steps_json = REPLACE(
    steps_json,
    '"type":"manual"',
    '"type":"approval"'
) WHERE steps_json LIKE '%"type":"manual"%';

-- Step 3: Rename "check" type to "run" and move "check" field to "command".
-- For check steps, the command was stored in the "check" field.
-- We rename the type and field in one pass.
UPDATE ops_runbooks SET steps_json = REPLACE(
    REPLACE(
        steps_json,
        '"type":"check"',
        '"type":"run"'
    ),
    '"check":',
    '"command":'
) WHERE steps_json LIKE '%"type":"check"%';

-- Step 4: Also migrate step types in existing step_results in run records.
UPDATE ops_runbook_runs SET step_results = REPLACE(
    REPLACE(
        REPLACE(
            step_results,
            '"type":"command"',
            '"type":"run"'
        ),
        '"type":"check"',
        '"type":"run"'
    ),
    '"type":"manual"',
    '"type":"approval"'
) WHERE step_results != '[]';
