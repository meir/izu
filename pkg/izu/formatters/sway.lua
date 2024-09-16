local formatter = {}
local izu = izu

izu.registerKeycode({
  ["super"] = "Mod4",
  "Shift",
  "Alt",
})

function formatter.keybind(parts)
  local bind = parts[1]
  local command = parts[2]
  return "bindsym " .. bind .. " exec " .. command
end

function formatter.command(parts)
  return table.concat(parts, "")
end

-- Super + { a, b } + XF68Media{Play,Pause}
-- ________________________________________
function formatter.base (parts)
  return table.concat(parts, "+")
end

-- Super + { a, b } + XF68Media{Play,Pause}
--         ^______^
function formatter.multiple (parts)
  return parts
end

-- Super + { a, b } + XF68Media{Play,Pause}
-- _____     _  _     _____________________
function formatter.single (parts)
  return table.concat(parts, "")
end

-- Super + { a, b } + XF68Media{Play,Pause}
--                             ^__________^
function formatter.single_part (parts)
  return parts
end

-- Super + { a, b } + XF68Media{Play,Pause}
-- ^^^^^     ^  ^     ^^^^^^^^^ ^^^^ ^^^^^
function formatter.string (part)
  return part
end

return formatter
