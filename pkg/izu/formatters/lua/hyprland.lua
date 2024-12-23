local formatter = {}
local izu = izu

-- modifier order for `bind = Super+Shift, exec, echo hellow world
local modifiers = {
  "Super",
  "Shift",
  "Alt",
  "Ctrl",
  "Shift_L",
  "Shift_R",
  "Alt_L",
  "Alt_R",
  "Ctrl_L",
  "Ctrl_R",
  "Super_L",
  "Super_R",
}

local function order_keys(bind)
  local mods = {}
  local keys = {}
  for _, v in ipairs(bind) do
    if izu.contains(modifiers, v) then
      table.insert(mods, v)
    else
      table.insert(keys, v)
    end
  end

  return {
    table.concat(mods, "+"),
    table.concat(keys, "+"),
  }
end

-- mousekeys for binds such as mouse:273, mouse:274, etc.
local mouse_keys = {
  ["mouse_lmb"] = "mouse:272",
  ["mouse_rmb"] = "mouse:273",
  ["mouse_mmb"] = "mouse:274",
}

local function replace_mousekey(key)
  if mouse_keys[key] then
    return mouse_keys[key]
  end

  if key:find("mouse_x") then
    local x = key:match("mouse_x(%d)")
    return "mouse:" .. (274 + tonumber(x))
  end

  return key
end

-- flags for binds, such as bindl, bindr, bindm, etc.
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

local function get_flags(flags)
  local bindflag = ""
  for _, v in pairs(flags) do
    if izu.contains(bindflags, v) then
      bindflag = bindflag .. v
    end
  end
  return bindflag
end

-- Formatter functions

function formatter.hotkey (args)
  local bindflag = get_flags(args.flags)
  return "bind" .. bindflag .. " = " .. table.concat(args.value, ", ")
end

function formatter.binding (args)
  if args.state == 1 then
    return table.concat(order_keys(args.value), ", ")
  end
  return table.concat(args.value, "")
end

function formatter.multiple (args)
  return args.value
end

function formatter.single (args)
  return replace_mousekey(table.concat(args.value, ""))
end

function formatter.string (args)
  return args.value
end

return formatter
