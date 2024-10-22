local formatter = {}
local izu = izu

local bindflags = {
  "l",
  "r",
  "e",
  "n",
  "m",
  "t",
  "i",
  "s",
  "d",
  "p",
}

local modifiers = {
  "super",
  "shift",
  "ctrl",
  "ctrl_l",
  "ctrl_r",
  "alt",
  "alt_l",
  "alt_r",
  "escape",
  "apostrophe",
}

local function flag(f)
  for _, v in pairs(bindflags) do
    if v == f then
      return v
    end
  end
  return ""
end

function formatter.hotkey (args)
  local flags = args.flags
  local bindflag = ""
  for _, v in pairs(flags) do
    bindflag = bindflag .. flag(v)
  end
  return "bind" .. bindflag .. " = " .. table.concat(args.value, ", ")
end

function formatter.binding (args)
  if args.state == 1 then
    return table.concat(args.value, ", ")
  end
  return table.concat(args.value, "")
end

function formatter.multiple (args)
  return args.value
end

function formatter.single (args)
  return table.concat(args.value, "")
end

function formatter.string (args)
  return args.value
end

return formatter
