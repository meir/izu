local formatter = {}
local izu = izu

function formatter.hotkey (args)
  return table.concat(args.value, "\n  ")
end

function formatter.binding (args)
  if args.state == 1 then
    return table.concat(args.value, " + ")
  end
  return table.concat(args.value, "")
end

function formatter.multiple (args)
  return "{" .. table.concat(args.value, ",") .. "}"
end

function formatter.single (args)
  return table.concat(args.value, "")
end

function formatter.string (args)
  return args.value
end

return formatter
