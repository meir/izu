local formatter = {}
local izu = izu


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

izu.registerKeycode({
  "SUPER",
  "SHIFT",
  "ALT",
})

function formatter.keybind(parts)
  local bind = parts[1]
  local command = parts[2]
  return "bind = " .. bind .. ", exec, " .. command
end

function formatter.command(parts)
  return table.concat(parts, "")
end

-- Super + { a, b } + XF68Media{Play,Pause}
-- ________________________________________
function formatter.base (parts)
  local modifier_list = {}
  local key_list = {}

  for _, part in ipairs(parts) do
    if izu.hasKey(modifiers, izu.lowercase(part)) then
      table.insert(modifier_list, part)
    else
      table.insert(key_list, part)
    end
  end

  return table.concat(modifier_list, "_") .. ", " .. table.concat(key_list, "&")
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
function formatter.string (part, section)
  return part
end

return formatter
